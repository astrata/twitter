/*
	Written by José Carlos Nieto <xiam@menteslibres.org>
	(c) 2012 Astrata Software http://astrata.mx

	MIT License
*/

package main

import (
	"flag"
	"fmt"
	"github.com/astrata/twitter"
	"github.com/garyburd/go-oauth/oauth"
	"log"
)

var consumerKey = flag.String("key", "", "Consumer key")
var consumerSecret = flag.String("secret", "", "Consumer secret")

func main() {

	flag.Parse()

	if *consumerKey == "" || *consumerSecret == "" {
		fmt.Printf("Get your Twitter token.\n\n")

		fmt.Printf("1. Register your app at http://dev.twitter.com.\n")
		fmt.Printf("2. Run this program with -key $CONSUMER_KEY -secret $CONSUMER_SECRET.\n\n")

		fmt.Printf("Arguments:\n")
		flag.PrintDefaults()
		return
	}

	client := twitter.New(oauth.Credentials{
		*consumerKey,
		*consumerSecret,
	})

	err := client.Setup()

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
}
