package splitter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Splitter struct {
	Document     string
	URI          string
	ResourceType string
}

func NewSplitter(uri string, document string, resourceType string) *Splitter {
	return &Splitter{
		Document:     document,
		ResourceType: resourceType,
		URI:          uri,
	}
}

func (s *Splitter) getDocument() ([]string, error) {
	var regex string
	switch s.ResourceType {
	case "js":
		regex = "<script[a-z|\\s|A-Z|0-9|.|\"|\\/|+|=]+src=[\"'].*?[\"'].*?><\\/script>"
	case "css":
		regex = "<link.*?href=[\"'].+?css.+?[\"'].*?[>\\/>]"
	default:
		return []string{}, errors.New("resource type doesn't exist")
	}
	reg := regexp.MustCompile(regex)
	splitted := reg.FindAllString(s.Document, -1)
	if len(splitted) == 0 {
		return []string{}, errors.New("document doesn't has " + s.ResourceType + " resource")
	}
	return splitted, nil
}

func (s *Splitter) GetResourceLink() ([]string, error) {
	var resources []string
	doc, err := s.getDocument()
	if err != nil {
		return []string{}, err
	}
	for _, d := range doc {
		var regex string
		switch s.ResourceType {
		case "js":
			regex = "src=[\"'].*?[\"']"
		case "css":
			regex = "href=[\"'].+?css.+?[\"']"
		default:
			return []string{}, errors.New("resource type doesn't exist")
		}
		m1, err := s.regexMatch(d, regex)
		if err != nil {
			fmt.Println("skipping", d)
			continue
		}
		res, err := s.regexMatch(m1, "['|\"](.*).*?['|\"]")
		if err != nil {
			fmt.Println("skipping " + d)
			continue
		}

		normalizedRes := s.regexSplit(res, "(\"|')")[1]
		_, err = s.regexMatch(normalizedRes, "^(https|http):\\/\\/[a-z|A-Z|.|0-9]+.*?\\/")
		if err != nil {
			fHost := strings.Split(normalizedRes, "/")
			uri := s.regexSplit(s.URI, "\\/[a-z|A-Z|.|\\-|0-9]+.(html|php)+?$")[0]
			if fHost[0] == "" && fHost[1] == "" {
				resources = append(resources, "http:"+normalizedRes)
			} else if fHost[0] == "" {
				resources = append(resources, uri+normalizedRes)
			} else {
				resources = append(resources, uri+"/"+normalizedRes)
			}
		} else {
			resources = append(resources, normalizedRes)
		}
	}
	if len(resources) == 0 {
		return []string{}, errors.New("document doesn't has a valid " + s.ResourceType + " url")
	}
	return resources, nil
}

func (s *Splitter) regexMatch(resource string, regex string) (string, error) {
	reg := regexp.MustCompile(regex)
	splitted := reg.FindString(resource)
	if splitted == "" {
		return "", errors.New("expression doesn't match any data")
	}
	return splitted, nil
}

func (s *Splitter) regexSplit(resource string, regex string) []string {
	reg := regexp.MustCompile(regex)
	splitted := reg.Split(resource, -1)
	return splitted
}
