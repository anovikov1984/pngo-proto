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

	FirstWay()
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
}

// only for publish
func ThirdWay() {
	ok := make(chan interface{})
	err := make(chan error)

	ch := pubnub.Publish().Channel("news").Success(ok).Error(err).PnChannel()

	ch <- 2
	ch <- 3
}
