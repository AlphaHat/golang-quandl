golang-quandl
=============

Quandl package for Go (golang)

Usage (anonymous):
```
q, _ := quandl.GetData("DMDRN/MSFT_MKT_CAP", "2000-01-01", "2013-01-05")

dates, marketCap := q.GetTimeSeriesData()
```

Usage (with auth token):
```
quandl.SetAuthToken("your auth token here")

q, _ := GetAllHistory("DMDRN/MSFT_MKT_CAP")

dates, marketCap := q.GetTimeSeriesData()

```

Running a search:
```
	body, err := Search("Apple Inc Short Interest")

	if err == nil {
		fmt.Printf("%s", body)
	}
```

Get all stocks in Quandl:
```
identifier, description := quandl.GetStockList()
```

Get all ETFs in Quandl:
```
identifier, description := quandl.GetETFList()
```

To use these identifiers you can simply use:
```
q, _ := quandl.GetData(identifier[0], "2013-01-01", "2013-01-05")

dates, values := q.GetTimeSeriesData()
```

The QuandResponse struct that is returned contains some metadata that may be useful for identifying attributes about the data that is returned.
```
type QuandlResponse struct {
	SourceCode string      `json:"source_code"`
	SourceName string      `json:"source_name"`
	Code       string      `json:"code"`
	Frequency  string      `json:"frequency"`
	FromDate   string      `json:"from_date"`
	ToDate     string      `json:"to_date"`
	Columns    []string    `json:"column_names"`
	Data       interface{} `json:"data"`
}
```
