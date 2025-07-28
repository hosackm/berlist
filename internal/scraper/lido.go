package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/hosackm/berlist/internal/event"
)

const TrimSet = `" „`

type LidoScraper struct {
	RootURL string
	Router  Router
}

func NewLidoScraper() *LidoScraper {
	scraper := &LidoScraper{
		RootURL: "https://www.lido-berlin.de",
		Router:  make(Router),
	}
	return scraper
}

func (l LidoScraper) Name() string {
	return "Lido Berlin"
}
func (l LidoScraper) AllowedDomains() []string {
	return []string{"www.lido-berlin.de", "lido-berlin.de"}
}
func (l LidoScraper) URL() string {
	return l.RootURL
}

func (l LidoScraper) GetRouter() Router {
	router := Router{
		"article.event-ticket": func(e *colly.HTMLElement) (*event.Event, error) {
			content := e.DOM.Find("div.event-ticket__content").First()

			artist := content.Find("a").First().Text()
			href, _ := content.Find("a").First().Attr("href")
			subtitle := content.Find("div.event-ticket__content__subtitle")

			sbuild := strings.Builder{}
			e.DOM.Find("div.event-ticket__meta__date").Children().Each(func(i int, s *goquery.Selection) {
				if i >= 3 {
					return
				}
				sbuild.WriteString(fmt.Sprintf(" %s", s.Text()))
			})
			startText := strings.Trim(sbuild.String(), " \t\r\n")

			start, err := time.Parse("2 January 2006", startText)
			if err != nil {
				start = time.Now()
			}

			return &event.Event{
				Date:     start,
				Name:     artist,
				Venue:    l.Name(),
				Subtitle: strings.Trim(subtitle.Text(), TrimSet),
				Link:     e.Request.AbsoluteURL(href),
			}, nil
		},
	}
	return router
}
