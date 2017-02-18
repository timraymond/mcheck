package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	ds, us, err := parsePage(tok)
	if err != nil {
		log.Println("Error while parsing: err:", err)
		os.Exit(1)
	}
	for _, d := range *ds {
		d.LineProtocol("channelstats", os.Stdout)
	}
	for _, u := range *us {
		u.LineProtocol("channelstats", os.Stdout)
	}
}

// parsePage extracts cable modem signal information from the stats page
// exposed at 192.168.100.1
func parsePage(t *html.Tokenizer) (ds *chanstat.ChannelStats, us *chanstat.UpstreamChannels, err error) {
	ds, err = consumeDownstreamStatsTable(t)
	if err != nil {
		return ds, us, err
	}
	us, err = consumeUpstreamStatsTable(t)
	if err != nil {
		return ds, us, err
	}
	err = consumeCodewordStats(t, ds)
	if err != nil {
		return ds, us, err
	}
	return ds, us, err
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

func consumeUpstreamStatsTable(t *html.Tokenizer) (*chanstat.UpstreamChannels, error) {
	discardElem(t, "tr")

	rawChanIds, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats := make(chanstat.UpstreamChannels, len(rawChanIds)-1)
	stats.AssignID(rawChanIds)

	rawFreqs, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats.AssignFreqs(rawFreqs)

	rawRangingIDs, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats.AssignRangingIDs(rawRangingIDs)

	rawSymRate, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats.AssignSymRate(rawSymRate)

	rawPowerLevels, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats.AssignPowerLevels(rawPowerLevels)

	discardElem(t, "tr") // toss out modulations

	rawRangingStatuses, err := parseRow(t)
	if err != nil {
		return nil, err
	}

	stats.AssignRangingStatus(rawRangingStatuses)

	return &stats, nil
}

func consumeDownstreamStatsTable(t *html.Tokenizer) (*chanstat.ChannelStats, error) {
	discardElem(t, "tr")

	rawChanIds, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats := make(chanstat.ChannelStats, len(rawChanIds)-1)
	stats.AssignID(rawChanIds)

	rawFreqs, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats.AssignFrequency(rawFreqs)

	rawSNRs, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats.AssignSNR(rawSNRs)

	rawModulations, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats.AssignMods(rawModulations)

	rawLevels, err := parseRow(t)
	if err != nil {
		return nil, err
	}
	stats.AssignLevels(rawLevels)

	return &stats, nil
}

func consumeCodewordStats(t *html.Tokenizer, ds *chanstat.ChannelStats) error {
	discardElem(t, "tr")
	discardElem(t, "tr")

	err := newParseRow(t, func(idx int, value string) error {
		unerrored, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
		if err != nil {
			return err
		}
		(*ds)[idx].CodewordStats.TotalUnerrored = unerrored
		return nil
	})

	if err != nil {
		return err
	}

	err = newParseRow(t, func(idx int, value string) error {
		correctable, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
		if err != nil {
			return err
		}
		(*ds)[idx].CodewordStats.TotalCorrectable = correctable
		return nil
	})

	if err != nil {
		return err
	}

	err = newParseRow(t, func(idx int, value string) error {
		uncorrectable, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
		if err != nil {
			return err
		}
		(*ds)[idx].CodewordStats.TotalUncorrectable = uncorrectable
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func newParseRow(t *html.Tokenizer, parser func(int, string) error) error {
	values, err := parseRow(t)
	if err != nil {
		return err
	}
	for idx, value := range values[1:] {
		if err := parser(idx, value); err != nil {
			return err
		}
	}
	return nil
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
