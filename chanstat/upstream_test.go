package chanstat_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/timraymond/mcheck/chanstat"
)

func Test_UpstreamAssignID(t *testing.T) {
	rawIDs := []string{"Channel ID", "3\u00a0 ", "1\u00a0 ", "2\u00a0 "}
	expected := chanstat.UpstreamChannels{
		chanstat.UpstreamChannel{ID: 3},
		chanstat.UpstreamChannel{ID: 1},
		chanstat.UpstreamChannel{ID: 2},
	}

	actual := make(chanstat.UpstreamChannels, len(rawIDs)-1)
	actual.AssignID(rawIDs)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Assigning upstream IDs failed. \nWant: \n\n %#v, \n Got: \n\n %#v", expected, actual)
	}
}

func Test_UpstreamAssignFreq(t *testing.T) {
	rawFreqs := []string{"Frequency", "23700000 Hz\u00a0", "36500000 Hz\u00a0", "30100000 Hz\u00a0"}
	expected := chanstat.UpstreamChannels{
		chanstat.UpstreamChannel{Frequency: uint64(23700000)},
		chanstat.UpstreamChannel{Frequency: uint64(36500000)},
		chanstat.UpstreamChannel{Frequency: uint64(30100000)},
	}

	actual := make(chanstat.UpstreamChannels, len(rawFreqs)-1)
	actual.AssignFreqs(rawFreqs)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Assigning upstream Frequencies failed. \nWant: \n\n %#v, \n Got: \n\n %#v", expected, actual)
	}
}

func Test_UpstreamAssignRangingIDs(t *testing.T) {
	rawRangingIDs := []string{"Ranging Service ID", "1761\u00a0", "1761\u00a0", "1761\u00a0"}
	expected := chanstat.UpstreamChannels{
		chanstat.UpstreamChannel{RangingServiceID: uint64(1761)},
		chanstat.UpstreamChannel{RangingServiceID: uint64(1761)},
		chanstat.UpstreamChannel{RangingServiceID: uint64(1761)},
	}

	actual := make(chanstat.UpstreamChannels, len(rawRangingIDs)-1)
	actual.AssignRangingIDs(rawRangingIDs)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Assigning ranging ids failed. \nWant: \n\n %#v, \n Got: \n\n %#v", expected, actual)
	}
}

func Test_UpstreamAssignSymRate(t *testing.T) {
	rawSymRates := []string{"Symbol Rate", "5.120 Msym/sec\u00a0", "5.120 Msym/sec\u00a0", "5.120 Msym/sec\u00a0"}
	expected := chanstat.UpstreamChannels{
		chanstat.UpstreamChannel{SymbolRate: float64(5.120)},
		chanstat.UpstreamChannel{SymbolRate: float64(5.120)},
		chanstat.UpstreamChannel{SymbolRate: float64(5.120)},
	}

	actual := make(chanstat.UpstreamChannels, len(rawSymRates)-1)
	actual.AssignSymRate(rawSymRates)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Assigning ranging ids failed. \nWant: \n\n %#v, \n Got: \n\n %#v", expected, actual)
	}
}

func Test_UpstreamAssignPowerLevels(t *testing.T) {
	rawPowerLevels := []string{"Power Level", "32 dBmV\u00a0", "31 dBmV\u00a0", "32 dBmV\u00a0", "33 dBmV\u00a0"}
	expected := chanstat.UpstreamChannels{
		chanstat.UpstreamChannel{PowerLevel: int64(32)},
		chanstat.UpstreamChannel{PowerLevel: int64(31)},
		chanstat.UpstreamChannel{PowerLevel: int64(32)},
		chanstat.UpstreamChannel{PowerLevel: int64(33)},
	}

	actual := make(chanstat.UpstreamChannels, len(rawPowerLevels)-1)
	actual.AssignPowerLevels(rawPowerLevels)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Assigning power levels failed. \nWant: \n\n %#v, \n Got: \n\n %#v", expected, actual)
	}
}

func Test_UpstreamAssignRanging(t *testing.T) {
	rawRangingStatus := []string{"Ranging Status ", "Success\u00a0", "Success\u00a0", "Success\u00a0", "Success\u00a0"}
	expected := chanstat.UpstreamChannels{
		chanstat.UpstreamChannel{RangingOK: true},
		chanstat.UpstreamChannel{RangingOK: true},
		chanstat.UpstreamChannel{RangingOK: true},
		chanstat.UpstreamChannel{RangingOK: true},
	}

	actual := make(chanstat.UpstreamChannels, len(rawRangingStatus)-1)
	actual.AssignRangingStatus(rawRangingStatus)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Assigning ranging statuses failed. \nWant: \n\n %#v, \n Got: \n\n %#v", expected, actual)
	}
}

func Test_Upstream_LineProtocol(t *testing.T) {
	bs := bytes.NewBufferString("")

	uchans := chanstat.UpstreamChannels{
		chanstat.UpstreamChannel{ID: 3, Frequency: uint64(23700000), PowerLevel: int64(41), SymbolRate: float64(4.815)},
	}

	for _, uc := range uchans {
		uc.LineProtocol("channel_stats", bs)
	}

	expected := "channel_stats,id=3,frequency=23700000,ranging_id=0 plevel=41,sym_rate=4.815\n"
	if actual := bs.String(); actual != expected {
		t.Error("Unsuccessful marshal to line protocol:\nwant:", expected, "\nGot:", actual)
	}
}
