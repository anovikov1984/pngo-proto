package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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
		resp, err := doRequest(url)
		if err != nil {
			e.errorChannel <- err
		}

		e.successChannel <- resp
	}()
}

func (e *Publish) buildUrl(msg interface{}) string {
	strMsg, _ := json.Marshal(msg)
	return "http://ps.pndsn.com/publish/pub-c-739aa0fc-3ed5-472b-af26-aca1b333ec52/sub-c-33f55052-190b-11e6-bfbc-02ee2ddab7fe/0/ch/0/%22" + string(strMsg) + "%22"
}

func (e *Publish) PnChannel() chan<- interface{} {
	err := e.validate()
	if err != nil {
		panic(err)
	}

	ch := make(chan interface{})

	go func() {
		for msg := range ch {

			url := e.buildUrl(msg)
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
	return http.Get(url)
}
