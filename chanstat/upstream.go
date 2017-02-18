package chanstat

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
)

type UpstreamChannel struct {
	ID               uint64
	Frequency        uint64
	RangingServiceID uint64
	SymbolRate       float64
	PowerLevel       int64
	Modulation       []Modulation
	RangingOK        bool
}

func (uc *UpstreamChannel) LineProtocol(measurement string, w io.Writer) {
	buf := bytes.NewBufferString(measurement)
	buf.WriteString(",")

	buf.WriteString("id=")
	buf.WriteString(strconv.Itoa(int(uc.ID)))
	buf.WriteString(",")

	buf.WriteString("frequency=")
	buf.WriteString(strconv.Itoa(int(uc.Frequency)))
	buf.WriteString(",")

	buf.WriteString("ranging_id=")
	buf.WriteString(strconv.Itoa(int(uc.RangingServiceID)))

	buf.WriteString(" ")

	buf.WriteString("plevel=")
	buf.WriteString(strconv.Itoa(int(uc.PowerLevel)))
	buf.WriteString(",")
	buf.WriteString("sym_rate=")
	buf.WriteString(strconv.FormatFloat(uc.SymbolRate, 'f', -1, 64))
	buf.WriteString(",")
	buf.WriteString("ranging_ok=")
	if uc.RangingOK {
		buf.WriteString("t")
	} else {
		buf.WriteString("f")
	}

	buf.WriteString("\n")
	buf.WriteTo(w)
}

type UpstreamChannels []UpstreamChannel

func (ucs *UpstreamChannels) AssignID(rawIDs []string) {
	for idx, rawID := range rawIDs[1:] {
		cleanID := strings.TrimSpace(rawID)
		ID, err := strconv.ParseUint(cleanID, 10, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*ucs)[idx].ID = ID
	}
}

func (ucs *UpstreamChannels) AssignFreqs(rawFreqs []string) {
	for idx, rawFreq := range rawFreqs[1:] {
		cleanFreq := strings.TrimSuffix(strings.TrimSpace(rawFreq), " Hz")
		freq, err := strconv.ParseUint(cleanFreq, 10, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*ucs)[idx].Frequency = freq
	}
}

func (ucs *UpstreamChannels) AssignRangingIDs(rawIDs []string) {
	for idx, rawID := range rawIDs[1:] {
		cleanID := strings.TrimSpace(rawID)
		id, err := strconv.ParseUint(cleanID, 10, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*ucs)[idx].RangingServiceID = id
	}
}

func (ucs *UpstreamChannels) AssignSymRate(rawSymRates []string) {
	for idx, rawSymRate := range rawSymRates[1:] {
		cleanSymRate := strings.TrimSuffix(strings.TrimSpace(rawSymRate), " Msym/sec")
		symRate, err := strconv.ParseFloat(cleanSymRate, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*ucs)[idx].SymbolRate = symRate
	}
}

func (ucs *UpstreamChannels) AssignPowerLevels(rawPLs []string) {
	for idx, rawPL := range rawPLs[1:] {
		cleanPL := strings.TrimSuffix(strings.TrimSpace(rawPL), " dBmV")
		pl, err := strconv.ParseInt(cleanPL, 10, 64)
		if err != nil {
			log.Println("err parsing: err", err)
		}
		(*ucs)[idx].PowerLevel = pl
	}
}

func (ucs *UpstreamChannels) AssignRangingStatus(rawRangingStatuses []string) {
	for idx, rawRng := range rawRangingStatuses[1:] {
		cleanRS := strings.TrimSpace(rawRng)
		if cleanRS == "Success" {
			(*ucs)[idx].RangingOK = true
		}
	}
}
