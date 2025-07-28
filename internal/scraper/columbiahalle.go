package scraper

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/hosackm/berlist/internal/event"
)

type ColumbiaHalleScraper struct {
	RootURL string
	Router  Router
}

func NewColumbiaHalleScraper() *ColumbiaHalleScraper {
	scraper := &ColumbiaHalleScraper{
		RootURL: "https://www.columbiahalle.berlin/events.html",
		Router:  make(Router),
	}
	return scraper
}

func (c ColumbiaHalleScraper) Name() string {
	return "ColumbiaHalle Berlin"
}

func (c ColumbiaHalleScraper) AllowedDomains() []string {
	return []string{"www.columbiahalle.berlin"}
}

func (c ColumbiaHalleScraper) URL() string {
	return c.RootURL
}

func (c ColumbiaHalleScraper) GetRouter() Router {
	router := Router{
		"div.event": func(e *colly.HTMLElement) (*event.Event, error) {
			dayStr := strings.TrimSpace(e.DOM.Find("p.event_datum_tag").Text())
			if dayStr == "" {
				return nil, nil
			}

			monthStr := ""
			e.DOM.PrevAllFiltered("div.eventlist_monat").EachWithBreak(func(_ int, s *goquery.Selection) bool {
				t := strings.TrimSpace(s.Text())
				if t == "" {
					return true
				}
				monthStr = t
				return false
			})

			dateStr := dayStr + " " + monthStr
			start, err := time.Parse("02 January 2006", dateStr)
			if err != nil {
				return nil, err
			}

			name := e.DOM.Find("div.event_kopf h2").Text()
			href, _ := e.DOM.Find("div.tickets div a").Attr("href")

			return &event.Event{
				Link:     e.Request.AbsoluteURL(href),
				Date:     start,
				Name:     name,
				Subtitle: "",
				Venue:    c.Name(),
			}, nil
		},
	}
	return router
}
