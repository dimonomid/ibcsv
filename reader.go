package ibcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

var (
	ErrNotEnoughFields  = errors.New("not enough fields")
	ErrWrongKind        = errors.New("wrong kind field")
	ErrUnexpectedHeader = errors.New("unexpected header")
	ErrMissingHeader    = errors.New("missing header")
)

type Table struct {
	// Name is the name of the table
	Name string

	// Fields is a slice of fields in a table, in the same order as they appear in CSV
	Fields []string

	// Rows is a slice of maps from the field name to the corresponding value.
	Rows []Row
}

type Row struct {
	Kind RowKind

	// Values is a map from the field name (as in Table.Fields slice) to the
	// corresponding value.
	Values map[string]string
}

// RowKind represents the row kind as found in Interactive Broker CSV files.
type RowKind string

const (
	RowKindData     RowKind = "data"
	RowKindSubtotal RowKind = "subtotal"
	RowKindTotal    RowKind = "total"
)

type Reader struct {
	csvReader *csv.Reader

	lastRecs []string
}

func NewReader(r io.Reader) *Reader {
	csvReader := csv.NewReader(r)
	// IBKR csv files have variable number of fields (when more than one table is
	// stored there), so we don't check it.
	csvReader.FieldsPerRecord = -1

	return &Reader{
		csvReader: csvReader,
	}
}

func (r *Reader) Read() (table *Table, err error) {
	for {
		recs := r.lastRecs
		r.lastRecs = nil

		if recs == nil {
			var err error
			recs, err = r.csvReader.Read()
			if err != nil {
				if err == io.EOF {
					return table, io.EOF
				}

				return nil, err
			}
		}

		if table == nil {
			table = &Table{}
		}

		if len(recs) < 3 {
			// Every line must contain at least the table name, then the "Header" or "Data",
			// then at least one more data field; so having less than 3 fields means a bad data.
			return nil, ErrNotEnoughFields
		}

		tableName := recs[0]

		var isHeader bool
		var rowKind RowKind

		switch recs[1] {
		case "Header":
			isHeader = true
		case "Data":
			rowKind = RowKindData
		case "SubTotal":
			rowKind = RowKindSubtotal
		case "Total":
			rowKind = RowKindTotal
		default:
			return nil, fmt.Errorf("%w: %s", ErrWrongKind, recs[1])
		}

		if tableName == table.Name && isHeader {
			return nil, fmt.Errorf("%w: table %s", ErrUnexpectedHeader, tableName)
		}

		if tableName != table.Name {
			if !isHeader {
				return nil, fmt.Errorf("%w: table %s", ErrMissingHeader, tableName)
			}

			if table.Name != "" {
				// A new table has started while we have the previous one to return
				r.lastRecs = recs
				return table, nil
			}

			// Started parsing table: remember field names

			table.Name = tableName
			table.Fields = recs[2:]

			continue
		}

		// Parsing table row
		row := Row{
			Kind:   rowKind,
			Values: make(map[string]string, len(recs)-2),
		}
		for i, v := range recs[2:] {
			row.Values[table.Fields[i]] = v
		}

		table.Rows = append(table.Rows, row)
	}
}
