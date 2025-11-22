package es

import (
	"log"

	"github.com/electricbubble/gadb"
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
