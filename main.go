package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var domain string
var t string

func Help() {
	fmt.Println(`getDomains - version 1.0
	-d	domain example.com
	-t Choose Onlye On [tlds/subdomains/all]
	`)
	os.Exit(0)
}

func main() {
	flag.StringVar(&domain, "d", "", "Add one Domain")
	flag.StringVar(&t, "t", "", "choose -t [tlds/subdomains/all] only one")
	flag.Parse()
	if t == "tlds" || t == "subdomains" || t == "all" && domain != "" {
		ch := make(chan map[string]bool, 1)
		go Sonar(domain, t, ch)
		for v := range ch {
			for mp := range v {
				fmt.Println(mp)
			}
		}
	} else {
		Help()

	}
}

func replaceS(s string) string {
	re := regexp.MustCompile(`"|^[|]`)
	s = re.ReplaceAllString(s, "")
	return strings.Trim(s, "[]")
}

func Sonar(d, t string, c chan map[string]bool) {
	m := make(map[string]bool)
	u := fmt.Sprintf("https://sonar.omnisint.io/%s/%s", t, d) // tlds subdomains all
	g, _ := http.Get(u)
	b, _ := ioutil.ReadAll(g.Body)
	s := strings.Split(string(b), ",")
	for _, v := range s {
		_, ok := m[replaceS(v)]
		if !ok {
			// fmt.Println(replaceS(v))
			m[replaceS(v)] = true
		}
	}
	c <- m
	close(c)
}
