package main

import (
	"github.com/elliotspeck/timboslice"
)

func main() {
	tim := timboslice.NewTim()
	tim.Connect()
	select {}
}
