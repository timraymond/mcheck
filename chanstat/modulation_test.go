package chanstat_test

import (
	"testing"

	"github.com/timraymond/mcheck/chanstat"
)

func Test_ParseModulation(t *testing.T) {
	modTests := []struct {
		modstr   string
		expected chanstat.Modulation
		errs     bool
	}{
		{"QPSK", chanstat.QPSK, false},
		{"QAM256", chanstat.QAM256, false},
		{"256QAM", chanstat.QAM256, false},
		{"64QAM", chanstat.QAM64, false},
		{"garbage", chanstat.Unknown, true},
	}

	for _, mTest := range modTests {
		actual, err := chanstat.ParseModulation(mTest.modstr)
		if err != nil && !mTest.errs {
			t.Error("Unexpected err parsing,", mTest.modstr, ", err was:", err)
		} else {
			if actual != mTest.expected {
				t.Error("Unexpected modulation: want:", mTest.expected, "got:", actual)
			}
		}
	}
}
