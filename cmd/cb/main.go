package main

import (
	"log"
)

func main() {
	log.SetFlags(0)

	cmd := rootCmd()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
