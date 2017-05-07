package main

import (
	"fmt"
	"log"

	"./endpoints"
)

var (
	pubnub *PubNub
	pnconf *PNConfiguration
)

type PubNub struct {
	pnconfig *PNConfiguration
}

func (pn *PubNub) Publish() *endpoints.Publish {
	return &endpoints.Publish{}
}

func NewPubNub(pnconfig *PNConfiguration) *PubNub {
	return &PubNub{pnconfig}
}

type PNConfiguration struct {
}

func main() {
	pnconf = &PNConfiguration{}
	pubnub = NewPubNub(pnconf)
	// ctx := context.Background()

	// FirstWay()
	// SecondWay()
	ThirdWay()
}

// Sync() generates a synchronous endpoit call and returns both response and
// error as a result
func FirstWay() {
	ok, err := pubnub.Publish().Channel("foo").Message("bar").Sync()
	if err != nil {
		log.Fatalf("Oooops! %s", err)
	}

	fmt.Println("1st way result", ok)
}

func SecondWay() {
	ok := make(chan interface{})
	err := make(chan error)

	pubnub.Publish().Channel("news").Success(ok).Error(err).Async()
	printResult(1, ok, err)
}

// only for publish
func ThirdWay() {
	ok := make(chan interface{})
	err := make(chan error)

	ch := pubnub.Publish().Channel("news").Success(ok).Error(err).PnChannel()
	go printResult(2, ok, err)

	ch <- 2
	ch <- 3

}

func printResult(times int, ok chan interface{}, err chan error) {
	for i := 0; i < times; i++ {
		select {
		case res := <-ok:
			fmt.Println(res)
		case er := <-err:
			fmt.Println(er)
		}
	}
}
