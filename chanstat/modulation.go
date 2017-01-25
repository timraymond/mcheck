package chanstat

import "errors"

// Modulation describes the type of modulation found on a particular channel
type Modulation int

const (
	Unknown Modulation = iota
	QPSK
	QAM64
	QAM256
)

func ParseModulation(modstr string) (Modulation, error) {
	switch modstr {
	case "QPSK":
		return QPSK, nil
	case "QAM256", "256QAM":
		return QAM256, nil
	case "64QAM":
		return QAM64, nil
	default:
		return Unknown, errors.New("Unknown modulation")
	}
}

func (m Modulation) String() string {
	switch m {
	case QPSK:
		return "QPSK"
	case QAM64:
		return "QAM64"
	case QAM256:
		return "QAM256"
	default:
		return "Unknown"
	}
}
