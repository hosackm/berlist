package scraper

import (
	"github.com/gocolly/colly/v2"
	"github.com/hosackm/berlist/internal/event"
)

// type RouteFunction colly.HTMLCallback

type RouteFunction func(e *colly.HTMLElement) (*event.Event, error)
type Router map[string]RouteFunction

type Scraper interface {
	Name() string

	URL() string

	AllowedDomains() []string

	GetRouter() Router
}
