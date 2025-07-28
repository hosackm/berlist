package scraper

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/hosackm/berlist/internal/event"
)

type SO36Scraper struct {
	RootURL string
	Router  Router
}

func NewSO36Scraper() *SO36Scraper {
	scraper := &SO36Scraper{
		RootURL: "https://www.so36.com/tickets",
		Router:  make(Router),
	}
	return scraper
}

func (s SO36Scraper) Name() string {
	return "SO36 Berlin"
}

func (s SO36Scraper) AllowedDomains() []string {
	return []string{"www.so36.de", "www.so36.com"}
}

func (s SO36Scraper) URL() string {
	return s.RootURL
}

func (s SO36Scraper) GetRouter() Router {
	router := Router{
		"li a": func(e *colly.HTMLElement) (*event.Event, error) {
			href := e.Attr("href")
			if !strings.HasPrefix(href, "/produkte") {
				return nil, nil
			}

			content := e.Text
			parts := strings.Split(content, " ")

			// Tickets {name} in Berlin am DD.MM.YYYY
			// 0        ...   -4   -3   -2    -1
			// artist := strings.Join(parts[1:len(parts)-4], " ")
			start, err := time.Parse("02.01.2006", parts[len(parts)-1])
			if err != nil {
				return nil, err
			}

			parts = parts[1 : len(parts)-4]
			artist := strings.Join(parts, " ")
			artist = strings.ReplaceAll(artist, "  ", " ")

			return &event.Event{
				Link:     e.Request.AbsoluteURL(href),
				Date:     start,
				Name:     artist,
				Subtitle: "",
				Venue:    s.Name(),
			}, nil
		},
		//
		// This is the full page but it's not returned when scraping.
		// Leaving it here for now in case I figure out how to use it.
		"div.type-ticket": func(e *colly.HTMLElement) (*event.Event, error) {
			artist := e.DOM.Find("div.ttli--where h2").First().Text()
			href, _ := e.DOM.Find("div.ttli-where a").Last().Attr("href")
			subtitle := e.DOM.Find("div.ttli-where small").First().Text()
			startText := e.DOM.Find("div.ttli--when span.date").First().Text()

			start, err := time.Parse("02.01.2006", strings.Trim(startText, " \t\n\r"))
			if err != nil {
				return nil, err
			}

			return &event.Event{
				Date:     start,
				Name:     artist,
				Venue:    "SO36 Berlin",
				Subtitle: strings.Trim(subtitle, TrimSet),
				Link:     e.Request.AbsoluteURL(href),
			}, nil
		},
	}
	return router
}
