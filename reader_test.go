package ibcsv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	comment     string
	data        string
	wantResults []testReadRes
}

type testReadRes struct {
	table *Table
	err   string
}

func TestReader(t *testing.T) {
	testCases := []testCase{
		testCase{comment: "normal data, without totals or sutbotals", // {{{
			data: `Statement,Header,Field Name,Field Value
Statement,Data,BrokerName,Interactive Brokers Central Europe Zrt.
Statement,Data,BrokerAddress,"Madach Imre ut 13-14, Floor 5, Budapest, 1075, Hungary"
Statement,Data,Title,Activity Statement
Statement,Data,Period,"February 1, 2022 - February 28, 2022"
Statement,Data,WhenGenerated,"2022-02-15, 02:21:06 EDT"
Withholding Tax,Header,Currency,Date,Description,Amount,Code
Withholding Tax,Data,USD,2022-01-08,FOO(qweqwe) Cash Dividend USD 2.00 per Share - US Tax,-6,
Withholding Tax,Data,USD,2022-01-14,BAR(asdasd) Cash Dividend USD 1.188 per Share - CA Tax,-10.49,
Withholding Tax,Data,USD,2022-01-14,BAZ(zxczxc) Cash Dividend USD 3.23 per Share - US Tax,-200.83,
Withholding Tax,Data,Total,,,-218.32,
Dividends,Header,Currency,Date,Description,Amount,Code
Dividends,Data,USD,2022-01-08,FOO(qweqwe) Cash Dividend USD 2.00 per Share (Ordinary Dividend),60,
Dividends,Data,USD,2022-01-14,BAR(asdasd) Cash Dividend USD 1.188 per Share (Ordinary Dividend),69.98,
Dividends,Data,USD,2022-01-14,BAZ(zxczxc) Cash Dividend USD 3.23 per Share (Ordinary Dividend),2008.3,
Dividends,Data,Total,,,4543.28,
`,
			wantResults: []testReadRes{
				testReadRes{
					table: &Table{
						Name:   "Statement",
						Fields: []string{"Field Name", "Field Value"},
						Rows: []Row{
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "BrokerName",
									"Field Value": "Interactive Brokers Central Europe Zrt.",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "BrokerAddress",
									"Field Value": "Madach Imre ut 13-14, Floor 5, Budapest, 1075, Hungary",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "Title",
									"Field Value": "Activity Statement",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "Period",
									"Field Value": "February 1, 2022 - February 28, 2022",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "WhenGenerated",
									"Field Value": "2022-02-15, 02:21:06 EDT",
								},
							},
						},
					},
				},
				testReadRes{
					table: &Table{
						Name:   "Withholding Tax",
						Fields: []string{"Currency", "Date", "Description", "Amount", "Code"},
						Rows: []Row{
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-08",
									"Description": "FOO(qweqwe) Cash Dividend USD 2.00 per Share - US Tax",
									"Amount":      "-6",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-14",
									"Description": "BAR(asdasd) Cash Dividend USD 1.188 per Share - CA Tax",
									"Amount":      "-10.49",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-14",
									"Description": "BAZ(zxczxc) Cash Dividend USD 3.23 per Share - US Tax",
									"Amount":      "-200.83",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "Total",
									"Date":        "",
									"Description": "",
									"Amount":      "-218.32",
									"Code":        "",
								},
							},
						},
					},
				},
				testReadRes{
					table: &Table{
						Name:   "Dividends",
						Fields: []string{"Currency", "Date", "Description", "Amount", "Code"},
						Rows: []Row{
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-08",
									"Description": "FOO(qweqwe) Cash Dividend USD 2.00 per Share (Ordinary Dividend)",
									"Amount":      "60",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-14",
									"Description": "BAR(asdasd) Cash Dividend USD 1.188 per Share (Ordinary Dividend)",
									"Amount":      "69.98",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-14",
									"Description": "BAZ(zxczxc) Cash Dividend USD 3.23 per Share (Ordinary Dividend)",
									"Amount":      "2008.3",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "Total",
									"Date":        "",
									"Description": "",
									"Amount":      "4543.28",
									"Code":        "",
								},
							},
						},
					},
				},
				testReadRes{
					err: "EOF",
				},
			},
		}, // }}}
		testCase{comment: "empty file", // {{{
			data: ``,
			wantResults: []testReadRes{
				testReadRes{
					table: nil,
					err:   "EOF",
				},
			},
		}, // }}}
		testCase{comment: "table with no data", // {{{
			data: `Statement,Header,Field Name,Field Value
Statement,Data,BrokerName,Interactive Brokers Central Europe Zrt.
Statement,Data,BrokerAddress,"Madach Imre ut 13-14, Floor 5, Budapest, 1075, Hungary"
Statement,Data,Title,Activity Statement
Statement,Data,Period,"February 1, 2022 - February 28, 2022"
Statement,Data,WhenGenerated,"2022-02-15, 02:21:06 EDT"
Withholding Tax,Header,Currency,Date,Description,Amount,Code
Dividends,Header,Currency,Date,Description,Amount,Code
Dividends,Data,USD,2022-01-08,FOO(qweqwe) Cash Dividend USD 2.00 per Share (Ordinary Dividend),60,
Dividends,Data,USD,2022-01-14,BAR(asdasd) Cash Dividend USD 1.188 per Share (Ordinary Dividend),69.98,
Dividends,Data,USD,2022-01-14,BAZ(zxczxc) Cash Dividend USD 3.23 per Share (Ordinary Dividend),2008.3,
Dividends,Data,Total,,,4543.28,
`,
			wantResults: []testReadRes{
				testReadRes{
					table: &Table{
						Name:   "Statement",
						Fields: []string{"Field Name", "Field Value"},
						Rows: []Row{
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "BrokerName",
									"Field Value": "Interactive Brokers Central Europe Zrt.",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "BrokerAddress",
									"Field Value": "Madach Imre ut 13-14, Floor 5, Budapest, 1075, Hungary",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "Title",
									"Field Value": "Activity Statement",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "Period",
									"Field Value": "February 1, 2022 - February 28, 2022",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Field Name":  "WhenGenerated",
									"Field Value": "2022-02-15, 02:21:06 EDT",
								},
							},
						},
					},
				},
				testReadRes{
					table: &Table{
						Name:   "Withholding Tax",
						Fields: []string{"Currency", "Date", "Description", "Amount", "Code"},
					},
				},
				testReadRes{
					table: &Table{
						Name:   "Dividends",
						Fields: []string{"Currency", "Date", "Description", "Amount", "Code"},
						Rows: []Row{
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-08",
									"Description": "FOO(qweqwe) Cash Dividend USD 2.00 per Share (Ordinary Dividend)",
									"Amount":      "60",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-14",
									"Description": "BAR(asdasd) Cash Dividend USD 1.188 per Share (Ordinary Dividend)",
									"Amount":      "69.98",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "USD",
									"Date":        "2022-01-14",
									"Description": "BAZ(zxczxc) Cash Dividend USD 3.23 per Share (Ordinary Dividend)",
									"Amount":      "2008.3",
									"Code":        "",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"Currency":    "Total",
									"Date":        "",
									"Description": "",
									"Amount":      "4543.28",
									"Code":        "",
								},
							},
						},
					},
				},
				testReadRes{
					err: "EOF",
				},
			},
		}, // }}}
		testCase{comment: "table with subtotal and total", // {{{
			data: `Trades,Header,One,Two,Three
Trades,Data,Foo,1,1
Trades,Data,Foo,3,1
Trades,SubTotal,Foo,4,2
Trades,Data,Bar,10,1
Trades,Data,Bar,30,1
Trades,SubTotal,Bar,40,2
Trades,Total,,44,4
`,
			wantResults: []testReadRes{
				testReadRes{
					table: &Table{
						Name:   "Trades",
						Fields: []string{"One", "Two", "Three"},
						Rows: []Row{
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"One":   "Foo",
									"Two":   "1",
									"Three": "1",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"One":   "Foo",
									"Two":   "3",
									"Three": "1",
								},
							},
							Row{
								Kind: RowKindSubtotal,
								Values: map[string]string{
									"One":   "Foo",
									"Two":   "4",
									"Three": "2",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"One":   "Bar",
									"Two":   "10",
									"Three": "1",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"One":   "Bar",
									"Two":   "30",
									"Three": "1",
								},
							},
							Row{
								Kind: RowKindSubtotal,
								Values: map[string]string{
									"One":   "Bar",
									"Two":   "40",
									"Three": "2",
								},
							},
							Row{
								Kind: RowKindTotal,
								Values: map[string]string{
									"One":   "",
									"Two":   "44",
									"Three": "4",
								},
							},
						},
					},
				},
				testReadRes{
					err: "EOF",
				},
			},
		}, // }}}
		testCase{comment: "repeating tables", // {{{
			data: `
Trades,Header,One,Two,Three
Trades,Data,Foo,1,1
Trades,Data,Foo,3,1
Trades,Total,Foo,4,2
Trades,Header,One,Two,Three
Trades,Data,Bar,10,1
Trades,Data,Bar,30,1
Trades,Total,Bar,40,2
`,
			wantResults: []testReadRes{
				testReadRes{
					table: &Table{
						Name:   "Trades",
						Fields: []string{"One", "Two", "Three"},
						Rows: []Row{
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"One":   "Foo",
									"Two":   "1",
									"Three": "1",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"One":   "Foo",
									"Two":   "3",
									"Three": "1",
								},
							},
							Row{
								Kind: RowKindTotal,
								Values: map[string]string{
									"One":   "Foo",
									"Two":   "4",
									"Three": "2",
								},
							},
						},
					},
				},
				testReadRes{
					table: &Table{
						Name:   "Trades",
						Fields: []string{"One", "Two", "Three"},
						Rows: []Row{
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"One":   "Bar",
									"Two":   "10",
									"Three": "1",
								},
							},
							Row{
								Kind: RowKindData,
								Values: map[string]string{
									"One":   "Bar",
									"Two":   "30",
									"Three": "1",
								},
							},
							Row{
								Kind: RowKindTotal,
								Values: map[string]string{
									"One":   "Bar",
									"Two":   "40",
									"Three": "2",
								},
							},
						},
					},
				},
				testReadRes{
					err: "EOF",
				},
			},
		}, // }}}
	}

	for i, tc := range testCases {
		buf := bytes.NewBuffer([]byte(tc.data))

		r := NewReader(buf)

		var gotResults []testReadRes

		for j := 0; ; j++ {
			if j >= 100 {
				assert.FailNow(t, "too many iterations")
			}

			table, err := r.Read()
			var errString string
			if err != nil {
				errString = err.Error()
			}
			gotResults = append(gotResults, testReadRes{
				table: table,
				err:   errString,
			})

			if err != nil {
				break
			}
		}

		assert.Equal(t, tc.wantResults, gotResults, "testCase #%d (%s)", i, tc.comment)
	}

}
