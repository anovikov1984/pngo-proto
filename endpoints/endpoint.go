package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Endpoint interface {
	Success() *Endpoint
	Error() *Endpoint
}

type TransactionalEndpoint interface {
	Sync() interface{}
	Async()
}

type BaseEndpoint struct {
	Endpoint
}

func (e *BaseEndpoint) Context(ctx context.Context) *BaseEndpoint {
	return e
}

type Publish struct {
	// BaseEndpoint
	Endpoint
	TransactionalEndpoint

	message string
	channel string

	successChannel chan interface{}
	errorChannel   chan error
}

func (e *Publish) Channel(channel string) *Publish {
	e.channel = channel
	return e
}

func (e *Publish) Message(msg string) *Publish {
	e.message = msg
	return e
}

func (e *Publish) Success(ch chan interface{}) *Publish {
	e.successChannel = ch
	return e
}

func (e *Publish) Error(ch chan error) *Publish {
	e.errorChannel = ch
	return e
}

func (e *Publish) validate() error {
	if e.successChannel == nil {
		return errors.New("pubnub: missing success go channel")
	}

	if e.errorChannel == nil {
		return errors.New("pubnub: missing error go channel")
	}

	return nil
}

func (e *Publish) Sync() (interface{}, error) {
	e.successChannel = make(chan interface{})
	e.errorChannel = make(chan error)

	e.Async()

	select {
	case resp := <-e.successChannel:
		return resp, nil
	case er := <-e.errorChannel:
		return nil, er
	}
}

func (e *Publish) Async() {
	err := e.validate()
	if err != nil {
		panic(err)
	}

	go func() {
		url := e.buildUrl(e.message)
		log.Println(url)
		resp, err := doRequest(url)
		if err != nil {
			e.errorChannel <- err
		}

		e.successChannel <- resp
	}()
}

func (e *Publish) buildUrl(msg interface{}) string {
	bytes, _ := json.Marshal(msg)
	strMsg := url.PathEscape(string(bytes))
	return "http://ps.pndsn.com/publish/pub-c-739aa0fc-3ed5-472b-af26-aca1b333ec52/sub-c-33f55052-190b-11e6-bfbc-02ee2ddab7fe/0/ch/0/" + string(strMsg) + "?uuid=gotest"
}

func (e *Publish) PnChannel() chan<- interface{} {
	err := e.validate()
	if err != nil {
		panic(err)
	}

	ch := make(chan interface{})

	go func() {
		log.Println("Waiting for channels")
		for msg := range ch {

			url := e.buildUrl(msg)
			log.Println(url)
			resp, err := doRequest(url)
			if err != nil {
				e.errorChannel <- err
			}

			e.successChannel <- resp
		}
	}()

	return ch
}

func doRequest(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println(string(body))
	return string(body), nil
}
