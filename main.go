package main

import (
	"github.com/LeeWaiHo/workflows/api/cmd"
	"log"
)

func main() {
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
