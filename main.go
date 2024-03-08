package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mahmoud-shabban/google-translate/cli"
)

var sourceLang string
var targetLang string
var sourceText string

var wg sync.WaitGroup

func init() {
	flag.StringVar(&sourceLang, "s", "en", "the source text language[en]")
	flag.StringVar(&targetLang, "t", "fr", "Target Language[fr]")
	flag.StringVar(&sourceText, "st", "", "Text to be Translated")
}
func main() {

	flag.Parse()
	// if with initialization statement
	if nf := flag.NFlag(); nf != 3 {
		fmt.Fprintf(os.Stderr, "Options\n")
		flag.PrintDefaults()
		// fmt.Printf("# flags: %d\n", nf)
		os.Exit(1)
	}

	strChan := make(chan string)

	wg.Add(1)

	reqBody := &cli.RequestBody{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		SourceText: sourceText,
	}

	go cli.RequestTranslate(reqBody, strChan, &wg)

	processedStr := strings.ReplaceAll(<-strChan, "+", " ")
	fmt.Fprintf(os.Stdout, "Translated Text:\n%s\n", processedStr)
	close(strChan)

	wg.Wait()
}
