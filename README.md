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

PACKAGE DOCUMENTATION
=====================
```

package quandl
    import "github.com/HedgeChart/golang-quandl/quandl"



FUNCTIONS

func GetAllSecurityList() ([]string, []string)
    GetAllSecurityList gets all the security identifiers and descriptions

func GetBitcoinList() ([]string, []string)
    GetBitcoinList just returns
    http://www.quandl.com/api/v1/datasets/BITCOIN/BITSTAMPUSD

func GetCommoditiesList() ([]string, []string)
    GetCommoditiesList gets all the Quandl codes and commodity descriptions

func GetDowConstituents() ([]string, []string)
    GetDowConstituents

func GetETFList() ([]string, []string)
    GetETFList gets all the Quandl codes and ETF descriptions

func GetETFTickerList() ([]string, []string)
    GetETFTickerList gets all the Quandl codes and ETF tickers

func GetEconomicDataList() ([]string, []string)
    Economic data (doesn't pertain to a particular security)
    GetEconomicDataList

func GetFTSE100Constituents() ([]string, []string)
    GetFTSE100Constituents

func GetFinancialRatiosList() ([]string, []string)
    GetFinancialRatiosList returns the list of Damordoran financial ratios.
    Currently this list is hard-coded into the file because Quandl does not
    provide a file from where to read. A caveat about using these is that
    you need to append the ticker in a particular way to use these ratios

func GetNasdaq100Constituents() ([]string, []string)
    GetNasdaq100Constituents

func GetNasdaqCompositeConstituents() ([]string, []string)
    GetNasdaqCompositeConstituents

func GetSP500Constituents() ([]string, []string)
    Index membership GetSP500Constituents

func GetSP500SectorMappings() ([]string, []string)
    Sector mappings GetSP500SectorMappings

func GetStockIndexList() ([]string, []string)
    GetStockIndexList gets all the Quandl codes and stock index descriptions

func GetStockList() ([]string, []string)
    GetStockList gets all the Quandl stock codes and descriptions

func GetStockTickerList() ([]string, []string)
    GetStockTickerList gets all the Quandl codes and tickers

func Search(query string) ([]byte, error)
    Search executes a query against the Quandl API and returns the JSON
    object as a byte stream. In future releases of this Go (golang) Quandl
    package this will return a native object instead of the json

func SetAuthToken(token string)
    SetAuthToken sets the auth token globally so that all subsequent calls
    that retrieve data from the Quandl API will use the auth token.


TYPES

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


func GetAllHistory(identifier string) (*QuandlResponse, error)
    GetAllHistory is similar to GetData except that it does not restrict a
    date range


func GetData(identifier string, startDate string, endDate string) (*QuandlResponse, error)
    GetData gets Quandl data for a particular identifier and a date range.
    You can optionally set the auth token before running this function so
    that you can make unlimited API calls instead of being limited to
    500/day.


func (q *QuandlResponse) GetTimeSeries(column string) ([]string, []float64)
    GetTimeSeries returns a date vector and the value vector for a
    particular column in the QuandlResponse

func (q *QuandlResponse) GetTimeSeriesColumn(column string) []float64
    GetTimeSeriesColumn returns the data from the Quandl response for a
    particular column. For some series, particularly stock data, multiple
    columns are returned. Using this method you can specify the specific
    column to extract.

func (q *QuandlResponse) GetTimeSeriesData() ([]float64, string)
    GetTimeSeriesData returns the most relevant data column from the Quandl
    response. In many cases you will not necessarily know beforehand what
    type of data is being requested and therefore cannot determine if it's
    stock data vs. economic data. In such cases, you can use this function
    to grab the column that is most likely relevant. The method also returns
    the most relevant

func (q *QuandlResponse) GetTimeSeriesDate() []string
    GetTimeSeriesDate returns the series of dates in the time series


```
