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
	Url string
	Count int
	Child []Knot
}

func getUrlLinks(url string) (correct []string){
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	list := regexp.MustCompile("href=\"[:.@/a-zA-Z0-9\\-_]+\"").FindAllString(string(body), -1)
	for _, link := range list {
		if regexp.MustCompile(Https+"?").MatchString(link) {
			if !regexp.MustCompile(ProtocolMatch + Host).MatchString(link) {
				continue
			}
		}
		link = regexp.MustCompile("("+Https+"?://"+Host+"|href=|\")").ReplaceAllString(link, "")
		if regexp.MustCompile("\\.(css|js|ico)").MatchString(link) {
			continue
		}
		correct = append(correct, link)
	}
	return
}

func main() {

	fmt.Println(getUrlLinks(Https + "://" + Host))

}
