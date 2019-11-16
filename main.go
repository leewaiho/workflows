package main

import (
	"github.com/LeeWaiHo/workflows/pkg/cmd"
	"log"
)

func main() {
	if e := cmd.Execute(); e != nil {
		log.Fatal(e)
	}
}
