package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/hosackm/berlist/internal/event"
)

func GetScrapers() []Scraper {
	return []Scraper{
		NewLidoScraper(),
		NewSO36Scraper(),
		NewColumbiaHalleScraper(),
	}
}

func RunScrapers() ([]*event.Event, error) {
	events := []*event.Event{}
	for _, s := range GetScrapers() {
		c := colly.NewCollector(colly.AllowedDomains(s.AllowedDomains()...))

		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Accept-Language", "en")
		})

		for selector, f := range s.GetRouter() {
			c.OnHTML(selector, func(e *colly.HTMLElement) {
				event, _ := f(e)
				if event != nil {
					events = append(events, event)
				}

			})
		}

		err := c.Visit(s.URL())
		if err != nil {
			return nil, err
		}

	}
	return events, nil
}
