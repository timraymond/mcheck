package chanstat

import (
	"log"
	"strconv"
	"strings"
)

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
