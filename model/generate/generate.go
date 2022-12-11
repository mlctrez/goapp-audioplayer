package main

import (
	"github.com/mlctrez/goapp-audioplayer/model/generate/generator"
	"log"
)

func main() {
	err := generator.Run()
	if err != nil {
		log.Fatal(err)
	}
}
