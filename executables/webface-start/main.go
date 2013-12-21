package main

import (
	"github.com/tsukanov/steaminfo-go/webface"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-dev" {
		webface.StartDev()
	} else {
		webface.Start()
	}
}
