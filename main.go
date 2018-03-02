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

type Knot struct {
    Url   string `json:"url"`
    Count int    `json:"count"`
    Child []Knot `json:"child"`
}

var list = createKnot("/")

func createKnot(url string) *Knot {
    return &Knot{url, 0, []Knot{}}
}

func (parent *Knot) addCountKnot() {
    parent.Count++
}

func (parent *Knot) search(url string) bool {
    if parent.Url == url {
        parent.addCountKnot()
        return true
    }
    for key, child := range parent.Child {
        if child.Url == url {
            child.addCountKnot()
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

func (parent *Knot) addKnot(child Knot) {
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

func stepInitListKnot(parent *Knot, url string) {
    urlList, err := getUrlLinks(Https + "://" + Host + url)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        for _, url := range urlList {
            state := parent.search(url)
            if state == false {
                parent.addKnot(Knot{url, 0, []Knot{}})
                stepInitListKnot(&list.Child[len(list.Child)-1], url)
            }
        }
    }

}

func main() {
    stepInitListKnot(list, "")
    data , _ := json.Marshal(list)
    fmt.Println(string(data))

}
