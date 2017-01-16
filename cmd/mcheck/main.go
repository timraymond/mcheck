package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

	consumeStatsTable(t, w, "downstream")
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

type ChannelStats []ChannelStat

func (cs *ChannelStats) AssignID(rawIDs []string) {
	for idx, rawID := range rawIDs[1:] {
		cleanID := strings.TrimSpace(rawID)
		ID, err := strconv.ParseUint(cleanID, 10, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*cs)[idx].ID = ID
	}
}

func (cs *ChannelStats) AssignFrequency(rawFreqs []string) {
	for idx, rawFreq := range rawFreqs[1:] {
		cleanFreq := strings.TrimSuffix(strings.TrimSpace(rawFreq), " Hz")
		freq, err := strconv.ParseUint(cleanFreq, 10, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*cs)[idx].Frequency = freq
	}
}

func (cs *ChannelStats) AssignSNR(rawSNRs []string) {
	for idx, rawSNR := range rawSNRs[1:] {
		cleanSNR := strings.TrimSuffix(strings.TrimSpace(rawSNR), " dB")
		snr, err := strconv.ParseInt(cleanSNR, 10, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*cs)[idx].SNR = snr
	}
}

func (cs *ChannelStats) AssignMods(rawMods []string) {
	for idx, rawMod := range rawMods[1:] {
		mod := strings.TrimSpace(rawMod)
		(*cs)[idx].Modulation = mod
	}
}

func (cs *ChannelStats) AssignLevels(rawLevels []string) {
	for idx, rawLevel := range rawLevels[1:] {
		cleanLevel := strings.TrimSuffix(strings.TrimSpace(rawLevel), " dBmV")
		level, err := strconv.ParseInt(cleanLevel, 10, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*cs)[idx].PowerLevel = level
	}
}

type ChannelStat struct {
	ID         uint64
	Frequency  uint64
	SNR        int64
	Modulation string
	PowerLevel int64
}

func (c *ChannelStat) LineProtocol(measurement string, w io.Writer) {
	buf := bytes.NewBufferString(measurement)
	buf.WriteString(",")

	buf.WriteString("id=")
	buf.WriteString(strconv.Itoa(int(c.ID)))
	buf.WriteString(",")

	buf.WriteString("frequency=")
	buf.WriteString(strconv.Itoa(int(c.Frequency)))
	buf.WriteString(" ")

	buf.WriteString("snr=")
	buf.WriteString(strconv.Itoa(int(c.SNR)))
	buf.WriteString(",")
	buf.WriteString("mod=")
	buf.WriteString(c.Modulation)
	buf.WriteString(",")
	buf.WriteString("plevel=")
	buf.WriteString(strconv.Itoa(int(c.PowerLevel)))

	buf.WriteString("\n")
	buf.WriteTo(w)
}

func consumeStatsTable(t *html.Tokenizer, w io.Writer, dir string) {
	discardElem(t, "tr")

	rawChanIds, err := parseRow(t)
	if err != nil {
		log.Println("Error occurred parsing row: err", err)
		os.Exit(1)
	}
	stats := make(ChannelStats, len(rawChanIds)-1)
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
		stat.LineProtocol("channel_stats", os.Stdout)
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
