package inspect

import (
	"bytes"
	"fmt"

	"github.com/gyuho/linux-inspect/proc"

	"github.com/gyuho/dataframe"
	"github.com/olekukonko/tablewriter"
)

// NSEntry represents network statistics.
// Simplied from 'NetDev'.
type NSEntry struct {
	Interface string

	ReceiveBytes    string
	ReceivePackets  uint64
	TransmitBytes   string
	TransmitPackets uint64

	// extra fields for sorting
	ReceiveBytesNum  uint64
	TransmitBytesNum uint64
}

// GetNS lists all '/proc/net/dev' statistics.
func GetNS() ([]NSEntry, error) {
	ss, err := proc.GetNetDev()
	if err != nil {
		return nil, err
	}
	ds := make([]NSEntry, len(ss))
	for i := range ss {
		ds[i] = NSEntry{
			Interface: ss[i].Interface,

			ReceiveBytes:    ss[i].ReceiveBytesParsedBytes,
			ReceivePackets:  ss[i].ReceivePackets,
			TransmitBytes:   ss[i].TransmitBytesParsedBytes,
			TransmitPackets: ss[i].TransmitPackets,

			ReceiveBytesNum:  ss[i].ReceiveBytesBytesN,
			TransmitBytesNum: ss[i].TransmitBytesBytesN,
		}
	}
	return ds, nil
}

const columnsNSToShow = 5

var columnsNSEntry = []string{
	"INTERFACE",

	"RECEIVE-BYTES", "RECEIVE-PACKETS",
	"TRANSMIT-BYTES", "TRANSMIT-PACKETS",

	// extra for sorting
	"RECEIVE-BYTES-NUM",
	"TRANSMIT-BYTES-NUM",
}

// ConvertNS converts to rows.
func ConvertNS(nss ...NSEntry) (header []string, rows [][]string) {
	header = columnsNSEntry
	rows = make([][]string, len(nss))
	for i, elem := range nss {
		row := make([]string, len(columnsNSEntry))
		row[0] = elem.Interface

		row[1] = elem.ReceiveBytes
		row[2] = fmt.Sprintf("%d", elem.ReceivePackets)
		row[3] = elem.TransmitBytes
		row[4] = fmt.Sprintf("%d", elem.TransmitPackets)

		row[5] = fmt.Sprintf("%d", elem.ReceiveBytesNum)
		row[6] = fmt.Sprintf("%d", elem.TransmitBytesNum)

		rows[i] = row
	}
	dataframe.SortBy(
		rows,
		dataframe.Float64DescendingFunc(5), // ReceiveBytesNum
		dataframe.Float64DescendingFunc(6), // TransmitBytesNum
	).Sort(rows)

	return
}

// StringNS converts in print-friendly format.
func StringNS(header []string, rows [][]string, topLimit int) string {
	buf := new(bytes.Buffer)
	tw := tablewriter.NewWriter(buf)
	tw.SetHeader(header[:columnsNSToShow:columnsNSToShow])

	if topLimit > 0 && len(rows) > topLimit {
		rows = rows[:topLimit:topLimit]
	}

	for _, row := range rows {
		tw.Append(row[:columnsNSToShow:columnsNSToShow])
	}
	tw.SetAutoFormatHeaders(false)
	tw.SetAlignment(tablewriter.ALIGN_RIGHT)
	tw.Render()

	return buf.String()
}
