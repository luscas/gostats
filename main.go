package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"regexp"
)

type StreamCast struct {
	StreamTitle      string
	Bitrate          string
	CurrentListeners string
	PeakListeners    string
	StreamGenre      string
	CurrentSong      string
}

func main() {
	godotenv.Load()

	e := echo.New()

	e.GET("/", statsHandler)

	e.Start(":8080")
}

func statsHandler(c echo.Context) error {
	url := os.Getenv("STREAMING_URL")

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var Result StreamCast

	doc.Find("table[cellpadding=\"2\"] tr").Each(func(i int, s *goquery.Selection) {
		var word = s.Find("td").Text()

		switch i {
		case 1:
			re := regexp.MustCompile(`Stream is up \((.*?)\) with ([0-9]+) of ([0-9]+) listeners`)
			match := re.FindStringSubmatch(word)
			Result.Bitrate = match[1]
			Result.CurrentListeners = match[2]
		case 2:
			re := regexp.MustCompile(`Listener Peak: ([0-9]+)`)
			match := re.FindStringSubmatch(word)
			Result.PeakListeners = match[1]
		case 4:
			re := regexp.MustCompile(`Stream Name: (.*)`)
			match := re.FindStringSubmatch(word)
			Result.StreamGenre = match[1]
		case 5:
			re := regexp.MustCompile(`Stream Genre\(s\): (.*)`)
			match := re.FindStringSubmatch(word)
			Result.StreamTitle = match[1]
		case 7:
			re := regexp.MustCompile(`Playing Now: (.*)`)
			match := re.FindStringSubmatch(word)
			Result.CurrentSong = match[1]
		}
	})

	c.JSON(http.StatusOK, Result)

	return nil
}
