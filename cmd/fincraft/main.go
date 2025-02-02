package main

import (
	"log"

	"fincraft/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg.FirstParam)
	log.Println(cfg.SecondParam)
}
