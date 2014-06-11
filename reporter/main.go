/*
Usage: reporter command

Commands:
	run-server - Start FastCGI server at 127.0.0.1:9000
	run-server-dev - Start development server at localhost:8080
*/
package main

import (
	"fmt"
	"os"

	"github.com/tsukanov/steamhistory/reporter/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No argument! See documentation to find out how to use this application.")
		return
	}

	switch os.Args[1] {
	case "run-server":
		server.Start()
	case "run-server-dev":
		server.StartDev()
	default:
		fmt.Println("Unknown command! See documentation.")
	}
}
