package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"flag"
	"github.com/schollz/progressbar/v3"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var (
	latitude      *string
	longitude     *string
	numberOfPages *int
	output        *string
)

type Map struct {
	ctx      context.Context
	BASE_URL string
	zoom     string
}

type Place struct {
	Name string
	Link string
}

func init() {
	latitude = flag.String("latitude", "35.1344161", "Provide latitude")
	longitude = flag.String("longitude", "33.3107072", "Provide longitude")
	numberOfPages = flag.Int("numberOfPages", 3, "Provide the number of pages you want to crawl")
	output = flag.String("output", "places.json", "Provide the output file name")
}

func (t *Map) getPlaces() []Place {
	var nodes []*cdp.Node
	err := chromedp.Run(t.ctx,
		chromedp.Sleep(5*time.Second),
		chromedp.Nodes("a.hfpxzc", &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Fatal(err)
	}

	places := []Place{}
	for _, node := range nodes {
		placeName, _ := node.Attribute("aria-label")
		placeLink, _ := node.Attribute("href")
		places = append(places, Place{Name: placeName, Link: placeLink})
	}

	return places
}

func (t *Map) scrollDown() {
	for i := 0; i < 5; i++ {
		err := chromedp.Run(t.ctx,
			chromedp.Sleep(5*time.Second),
			chromedp.Evaluate(`document.querySelector("#QA0Szd > div > div > div.w6VYqd > div.bJzME.tTVLSc > div > div.e07Vkf.kA9KIf > div > div > div.m6QErb.DxyBCb.kA9KIf.dS8AEf.ecceSd > div.m6QErb.DxyBCb.kA9KIf.dS8AEf.ecceSd").scroll(0,document.querySelector("#QA0Szd > div > div > div.w6VYqd > div.bJzME.tTVLSc > div > div.e07Vkf.kA9KIf > div > div > div.m6QErb.DxyBCb.kA9KIf.dS8AEf.ecceSd > div.m6QErb.DxyBCb.kA9KIf.dS8AEf.ecceSd").scrollHeight)`, nil),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (t *Map) nextPage() {
	err := chromedp.Run(t.ctx,
		chromedp.Click("#ppdPk-Ej1Yeb-LgbsSe-tJiF1e", chromedp.ByQuery),
		chromedp.Sleep(4*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *Map) dumpToJson(places map[string]Place, outputFile string) error {
	data := []Place{}
	for _, place := range places {
		data = append(data, place)
	}
	output, _ := json.Marshal(data)
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(output)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	flag.Parse()

	var m Map
	ctx, cancel := chromedp.NewContext(context.Background())
	m.ctx = ctx
	m.BASE_URL = "https://www.google.com/maps/search/Restaurants/@"
	m.zoom = "13z"

	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(m.BASE_URL+*latitude+","+*longitude+"/"+m.zoom),
	)
	if err != nil {
		log.Fatal(err)
	}

	places := make(map[string]Place)

	bar := progressbar.Default(100)
	for i := 0; i < *numberOfPages; i++ {
		bar.Add(100 / *numberOfPages)
		m.scrollDown()
		for _, place := range m.getPlaces() {
			places[place.Name] = place
		}
		m.nextPage()
	}

	bar.Set(100)

	m.dumpToJson(places, *output)

}
