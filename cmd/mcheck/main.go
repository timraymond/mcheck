package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/timraymond/mcheck/chanstat"
	"golang.org/x/net/html"
)

func main() {
	resp, err := http.Get("http://192.168.100.1/cmSignalData.htm")
	if err != nil {
		fmt.Println("Unable to fetch modem stats: err:", err)
		os.Exit(1)
	}
	tok := html.NewTokenizer(resp.Body)
	parsePage(tok, os.Stdout)
}

// parsePage extracts cable modem signal information from the stats page
// exposed at 192.168.100.1
func parsePage(t *html.Tokenizer, w io.Writer) {
	//discardElem(t, "table") // discard page header
	//discardElem(t, "table") // discard navigation
	//discardElem(t, "table") // discard page information

	consumeDownstreamStatsTable(t, w, "downstream")
	consumeUpstreamStatsTable(t, w)
	//consumeStatsTable(t, w, "upstream")
	//consumeCodewordsTable(t, w, "codewords")
}

func discardElem(t *html.Tokenizer, elem string) {
	for {
		tt := t.Next()
		if tt == html.EndTagToken {
			tok := t.Token()
			if tok.Data == elem {
				break
			}
		}
	}
}

func consumeUpstreamStatsTable(t *html.Tokenizer, w io.Writer) {
	discardElem(t, "tr")

	rawChanIds, err := parseRow(t)
	if err != nil {
		log.Println("Error occurred parsing row: err", err)
		os.Exit(1)
	}
	stats := make(chanstat.UpstreamChannels, len(rawChanIds)-1)
	stats.AssignID(rawChanIds)

	rawFreqs, err := parseRow(t)
	if err != nil {
		log.Println("Error occurred parsing row: err", err)
		os.Exit(1)
	}
	stats.AssignFreqs(rawFreqs)

	rawRangingIDs, err := parseRow(t)
	if err != nil {
		log.Println("Error occurred parsing row: err", err)
		os.Exit(1)
	}
	stats.AssignRangingIDs(rawRangingIDs)

	rawSymRate, err := parseRow(t)
	if err != nil {
		log.Println("Error occurred parsing row: err", err)
		os.Exit(1)
	}
	stats.AssignSymRate(rawSymRate)

	rawPowerLevels, err := parseRow(t)
	if err != nil {
		log.Println("Error occurred parsing row: err", err)
		os.Exit(1)
	}
	stats.AssignPowerLevels(rawPowerLevels)

	discardElem(t, "tr") // toss out modulations

	rawRangingStatuses, err := parseRow(t)
	if err != nil {
		log.Println("Error occurred parsing row: err", err)
		os.Exit(1)
	}

	stats.AssignRangingStatus(rawRangingStatuses)

	for _, stat := range stats {
		stat.LineProtocol("channelstats", os.Stdout)
	}
}

func consumeDownstreamStatsTable(t *html.Tokenizer, w io.Writer, dir string) {
	discardElem(t, "tr")

	rawChanIds, err := parseRow(t)
	if err != nil {
		log.Println("Error occurred parsing row: err", err)
		os.Exit(1)
	}
	stats := make(chanstat.ChannelStats, len(rawChanIds)-1)
	stats.AssignID(rawChanIds)

	rawFreqs, err := parseRow(t)
	if err != nil {
		log.Println("Error parsing frequencies: err", err)
		os.Exit(1)
	}
	stats.AssignFrequency(rawFreqs)

	rawSNRs, err := parseRow(t)
	if err != nil {
		log.Println("Error parsing frequencies: err", err)
		os.Exit(1)
	}
	stats.AssignSNR(rawSNRs)

	rawModulations, err := parseRow(t)
	if err != nil {
		log.Println("Error parsing frequencies: err", err)
		os.Exit(1)
	}
	stats.AssignMods(rawModulations)

	rawLevels, err := parseRow(t)
	if err != nil {
		log.Println("Error parsing frequencies: err", err)
		os.Exit(1)
	}
	stats.AssignLevels(rawLevels)

	for _, stat := range stats {
		stat.LineProtocol("channelstats", os.Stdout)
	}
}

func parseRow(t *html.Tokenizer) ([]string, error) {
	out := []string{}
	consume := false
	innerTables := 0
	for {
		tt := t.Next()
		switch {
		case tt == html.ErrorToken:
			if err := t.Err(); err.Error() == "EOF" {
				return out, nil
			} else {
				return out, err
			}
		case tt == html.StartTagToken:
			tok := t.Token()
			consume = tok.Data == "td"
			if tok.Data == "table" {
				innerTables++
			}
		case tt == html.TextToken:
			if consume {
				tok := t.Token()
				out = append(out, tok.Data)
			}
		case tt == html.EndTagToken:
			tok := t.Token()
			switch tok.Data {
			case "td":
				consume = false
			case "tr":
				if innerTables == 0 {
					return out, nil
				}
			case "table":
				innerTables--
			}
		}
	}

	return out, nil
}
