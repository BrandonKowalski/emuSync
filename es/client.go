package es

import (
	"github.com/electricbubble/gadb"
	"log"
)

var client gadb.Client

type EmuSync struct{}

func init() {
	var err error
	client, err = gadb.NewClient()

	if err != nil {
		log.Fatal(err)
	}
}
