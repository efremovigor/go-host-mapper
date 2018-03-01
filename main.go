package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
)

const Http = "http"
const Https = Http + "s"
const ProtocolMatch = Https + "?://"
const Host = "gardenmoto.ru"

type Knot struct {
	Url   string
	Count int
	Child []Knot
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

func getUrlLinks(url string) (correct []string) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	list := regexp.MustCompile("href=\"[:.@/a-zA-Z0-9\\-_]+\"").FindAllString(string(body), -1)
	for _, link := range list {
		if regexp.MustCompile(Https + "?").MatchString(link) {
			if !regexp.MustCompile(ProtocolMatch + Host).MatchString(link) {
				continue
			}
		}
		link = regexp.MustCompile("(" + Https + "?://" + Host + "|href=|\")").ReplaceAllString(link, "")
		if regexp.MustCompile("\\.(css|js|ico)").MatchString(link) {
			continue
		}
		correct = append(correct, link)
	}
	return
}

func stepInitListKnot(parent *Knot, url string){
	for _, url := range getUrlLinks(Https + "://" + Host + url) {
		state := parent.search(url)
		if state == false {
			parent.addKnot(Knot{url, 0, []Knot{}})
			stepInitListKnot(&list.Child[len(list.Child)-1], url)
		}
	}
}

func main() {
	stepInitListKnot(list,"")
	fmt.Println(list)

}
