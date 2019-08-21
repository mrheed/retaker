package main

import (
	"flag"
	"fmt"
	req "github.com/syahidnurrohim/retaker/requester"
	spltr "github.com/syahidnurrohim/retaker/splitter"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	url := flag.String("s", "", "website destination")
	resType := flag.String("t", "", "resource type, current available (js, css)")
	path := flag.String("o", "", "output directory")
	flag.Parse()
	if *url == "" {
		fmt.Println("error: empty parameter -s")
		return
	}
	if *resType == "" {
		fmt.Println("error: empty parameter -t")
		return
	}
	if *path == "" {
		fmt.Println("error: empty parameter -o")
		return
	}
	requester := req.NewRequester(*url)
	body, err := requester.GetBody()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	splitter := spltr.NewSplitter(requester.URI, body, *resType)
	res, err := splitter.GetResourceLink()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	CreateDirIfNotExist(*path + "/")
	for _, d := range res {
		requester.URI = d
		body1, err := requester.GetBody()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		ketok := strings.Split(d, "/")
		filename := ketok[len(ketok)-1]
		d1 := []byte(body1)
		err = ioutil.WriteFile(*path+filename, d1, 0644)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(filename, "downloaded")
	}
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
