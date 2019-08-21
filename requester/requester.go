package requester

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Requester struct {
	URI string
}

func NewRequester(URI string) *Requester {
	return &Requester{
		URI: URI,
	}
}

func (r *Requester) GetBody() (string, error) {
	resp, err := http.Get(r.URI)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		splitted := strings.Split(r.URI, "/")
		return "", errors.New("unable to make request on " + splitted[len(splitted)-1] + " with returned status: " + resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
