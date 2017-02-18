package chanstat

import (
	"bytes"
	"io"
	"strconv"
)

type CodewordStats struct {
	TotalUnerrored     uint64
	TotalCorrectable   uint64
	TotalUncorrectable uint64
}

type ChannelStat struct {
	ID            uint64
	Frequency     uint64
	SNR           int64
	Modulation    string
	PowerLevel    int64
	CodewordStats CodewordStats
}

func (c *ChannelStat) LineProtocol(measurement string, w io.Writer) {
	buf := bytes.NewBufferString(measurement)
	buf.WriteString(",")

	buf.WriteString("id=")
	buf.WriteString(strconv.Itoa(int(c.ID)))
	buf.WriteString(",")
	buf.WriteString("direction=down,")

	buf.WriteString("frequency=")
	buf.WriteString(strconv.Itoa(int(c.Frequency)))
	buf.WriteString(" ")

	buf.WriteString("snr=")
	buf.WriteString(strconv.Itoa(int(c.SNR)))
	buf.WriteString(",")
	buf.WriteString("mod=\"")
	buf.WriteString(c.Modulation)
	buf.WriteString("\",")
	buf.WriteString("plevel=")
	buf.WriteString(strconv.Itoa(int(c.PowerLevel)))
	buf.WriteString(",")
	buf.WriteString("unerrored=")
	buf.WriteString(strconv.Itoa(int(c.CodewordStats.TotalUnerrored)))
	buf.WriteString(",")
	buf.WriteString("correctable=")
	buf.WriteString(strconv.Itoa(int(c.CodewordStats.TotalCorrectable)))
	buf.WriteString(",")
	buf.WriteString("uncorrectable=")
	buf.WriteString(strconv.Itoa(int(c.CodewordStats.TotalUncorrectable)))

	buf.WriteString("\n")
	buf.WriteTo(w)
}
