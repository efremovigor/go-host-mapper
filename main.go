package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "regexp"
    "errors"
    "encoding/json"
)

const Http = "http"
const Https = Http
const ProtocolMatch = Https + "?://"
const Host = "anton.shevchuk.name"

var badUrl []string

func addBadUrl(url string) {
    badUrl = append(badUrl, url)
}

type Unit struct {
    Url   string `json:"url"`
    Count int    `json:"count"`
    Child []Unit `json:"child"`
}

var list = createUnit("/")

func createUnit(url string) *Unit {
    return &Unit{url, 0, []Unit{}}
}

func (parent *Unit) addCountUnit() {
    parent.Count++
}

func (parent *Unit) search(url string) bool {
    if parent.Url == url {
        parent.addCountUnit()
        return true
    }
    for key, child := range parent.Child {
        if child.Url == url {
            child.addCountUnit()
            parent.Child[key] = child
            return true
        }
        state := child.search(url)
        if state == true {
            return true
        }
    }
    return false
}

func (parent *Unit) addUnit(child Unit) {
    parent.Child = append(parent.Child, child)
}

func getUrlLinks(url string) (correct []string, err error) {
    resp, err := http.Get(url)
    if err != nil {
        err = errors.New("ошибка запроса - " + err.Error())
        addBadUrl(url)
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        err = errors.New("ошибка чтения тела - " + err.Error())
        return
    }
    list := regexp.MustCompile("href=\"[:.@/a-zA-Z0-9\\-_]+\"").FindAllString(string(body), -1)
    for _, link := range list {
        if regexp.MustCompile(Https + "?").MatchString(link) {
            if !regexp.MustCompile(ProtocolMatch + Host).MatchString(link) {
                continue
            }
        }
        link = regexp.MustCompile("(" + Https + "?://" + Host + "|href=|\")").ReplaceAllString(link, "")
        if regexp.MustCompile("\\.(css|js|ico|jpg|png)").MatchString(link) {
            continue
        }
        correct = append(correct, link)
    }
    return
}

func stepInitListUnit(parent *Unit, url string) {
    urlList, err := getUrlLinks(Https + "://" + Host + url)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        for _, url := range urlList {
            state := parent.search(url)
            if state == false {
                parent.addUnit(Unit{url, 0, []Unit{}})
                stepInitListUnit(&list.Child[len(list.Child)-1], url)
            }
        }
    }

}

func main() {
    stepInitListUnit(list, "")
    data , _ := json.Marshal(list)
    fmt.Println(string(data))

}
