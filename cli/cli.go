package cli

import (
	"log"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	SourceText string
	TargetLang string
}

const TranslateUrl = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, strChan chan string, wg *sync.WaitGroup) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", TranslateUrl, nil)

	if err != nil {
		log.Fatalf("1. There was a Problem: %s\n", err)
	}

	query := req.URL.Query()

	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	// t means translation of source text
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("2. There was a Problem: %s\n", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		strChan <- "You have been rate limited, try again later."
		wg.Done()
		return
	}

	parsedJson, err := gabs.ParseJSONBuffer(resp.Body)

	if err != nil {
		log.Fatalf("3. There was a Problem: %s\n", err)
	}

	nestedOne, err := parsedJson.ArrayElement(0)

	if err != nil {
		log.Fatalf("4. There was a Problem: %s\n", err)
	}

	nestedTwo, err := nestedOne.ArrayElement(0)

	if err != nil {
		log.Fatalf("5. There was a Problem: %s\n", err)
	}

	translatedStr, err := nestedTwo.ArrayElement(0)

	if err != nil {
		log.Fatalf("6. There was a Problem: %s\n", err)
	}

	strChan <- translatedStr.Data().(string)

	wg.Done()

}
