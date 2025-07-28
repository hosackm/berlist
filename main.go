package main

import (
	"encoding/json"
	"os"

	"github.com/hosackm/berlist/internal/scraper"
)

func main() {
	events, err := scraper.RunScrapers()
	if err != nil {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]string{"error": err.Error()})
		os.Exit(1)
	}

	_ = json.NewEncoder(os.Stdout).Encode(events)
}
