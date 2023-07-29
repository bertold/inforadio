package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/playwright-community/playwright-go"
)

func main() {
	url, err := getFeedURL()
	if err != nil {
		log.Fatalf("failed to get feed URL: %v", err)
	}

	err = createFeed(*url)
	if err != nil {
		log.Fatalf("failed to create feed: %v", err)
	}
}

func getFeedURL() (*string, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start playwright: %v", err)
	}
	defer func() {
		err = pw.Stop()
		if err != nil {
			log.Printf("could not stop playwright: %v", err)
		}
	}()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}
	defer func() {
		err = browser.Close()
		if err != nil {
			log.Printf("failed to close browser: %v", err)
		}
	}()
	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %v", err)
	}

	if _, err = page.Goto("https://infostart.hu/inforadio/napinfo"); err != nil {
		return nil, fmt.Errorf("could not goto: %v", err)
	}

	// accept cookies
	err = page.Locator("button[mode='primary']").Click()
	if err != nil {
		return nil, fmt.Errorf("could not click: %v", err)
	}

	// click on the first play link
	entries, err := page.Locator("span[data-bs-toggle='modal']").All()
	if err != nil {
		return nil, fmt.Errorf("could not get entries: %v", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("could not find links")
	}
	err = entries[0].Click()
	if err != nil {
		return nil, fmt.Errorf("failed to click: %v", err)
	}

	// extract URL
	href, err := page.Locator("div.infoplayer-head-icons > a").GetAttribute("href")
	if err != nil {
		return nil, fmt.Errorf("failed to find link: %v", err)
	}

	url := strings.Replace(href, "https://chtbl.com/track/GB95AD/dts.podtrac.com/redirect.mp3", "https:/", 1)
	return &url, nil
}

/*
	{
	  "uid": "urn:uuid:${UUID}",
	  "updateDate": "${UPDATED}",
	  "titleText": "Latest news from Klub Radio",
	  "mainText": "",
	  "streamUrl": "${STREAM_URL}",
	  "redirectionUrl": "https://www.klubradio.hu/"
	}
*/
type AlexaFeed struct {
	UID            string `json:"uid"`
	UpdateDate     string `json:"updateDate"`
	TitleText      string `json:"titleText"`
	MainText       string `json:"mainText"`
	StreamUrl      string `json:"streamUrl"`
	RedirectionUrl string `json:"redirectionUrl"`
}

func createFeed(url string) error {
	feed := AlexaFeed{
		UID:            "urn:uuid:" + uuid.New().String(),
		UpdateDate:     time.Now().Format(time.RFC3339),
		TitleText:      "Latest news from Info Radio",
		MainText:       "",
		StreamUrl:      url,
		RedirectionUrl: "https://infostart.hu/inforadio/napinfo",
	}

	data, err := json.Marshal(feed)
	if err != nil {
		return err
	}

	return os.WriteFile("inforadio.json", data, 0644)
}
