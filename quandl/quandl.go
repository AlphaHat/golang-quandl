package quandl

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	quandlApiRoot    = "http://www.quandl.com/api/v1/datasets/"
	quandlSearchRoot = "http://www.quandl.com/api/v1/datasets.json"

	quandlStockList     = "https://s3.amazonaws.com/quandl-static-content/quandl-stock-code-list.csv"
	quandlStockWikiList = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/WIKI_tickers.csv"
	sectorList          = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/Stock+Exchanges/stockinfo.csv"
	etfList             = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/ETFs.csv"
	stockIndexList      = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/Stock+Exchanges/Indicies.csv"
	mutualFundList      = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/Stock+Exchanges/funds.csv"

	// Code already contains the source
	commoditiesList = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/commodities.csv"

	// Source is in the file and needs to be pre-pended
	currencyList = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/currencies.csv"

	spxConstituents             = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/Indicies/SP500.csv"
	dowConstituents             = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/Indicies/dowjonesIA.csv"
	nasdaqCompositeConstituents = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/Indicies/NASDAQComposite.csv"
	nasdaq100Constituents       = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/Indicies/nasdaq100.csv"
	ftse100Constituents         = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/Indicies/FTSE100.csv"

	// Note that this data is pipe-delimited rather than comma delimited
	economicData = "https://s3.amazonaws.com/quandl-static-content/Ticker+CSV%27s/FRED/fred_allcodes.csv"

	format = ".json"
)

var authToken string

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

/*type TimeSeriesDataPoint struct {
	Date string
	Data float64
}

type TimeSeries struct {
	Data []TimeSeriesDataPoint
}*/

// SetAuthToken sets the auth token globally so that all subsequent calls that
// retrieve data from the Quandl API will use the auth token.
func SetAuthToken(token string) {
	authToken = token
}

func assembleQueryURL(query string) string {
	var url string

	if authToken == "" {
		fmt.Printf("No auth token set. API calls are limited.\n")
		url = fmt.Sprintf("%s?query=%s", quandlSearchRoot, query)
	} else {
		url = fmt.Sprintf("%s?query=%s&auth_token=%s", quandlSearchRoot, query, authToken)
	}
	return url
}

func assembleURLwithDates(identifier string, startDate string, endDate string) string {
	var url string

	if authToken == "" {
		fmt.Printf("No auth token set. API calls are limited.\n")
		url = fmt.Sprintf("%s%s%s?trim_start=%s&trim_end=%s", quandlApiRoot, identifier, format, startDate, endDate)
	} else {
		url = fmt.Sprintf("%s%s%s?trim_start=%s&trim_end=%s&auth_token=%s", quandlApiRoot, identifier, format, startDate, endDate, authToken)
	}
	return url
}

func assembleURLwithoutDates(identifier string) string {
	var url string
	if authToken == "" {
		fmt.Printf("No auth token set. API calls are limited.\n")
		url = fmt.Sprintf("%s%s%s", quandlApiRoot, identifier, format)
	} else {
		url = fmt.Sprintf("%s%s%s?auth_token=%s", quandlApiRoot, identifier, format, authToken)
	}
	return url
}

func readBytesFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func getDataFromURL(url string) (*QuandlResponse, error) {
	//fmt.Printf("%s\n", url)

	body, err := readBytesFromUrl(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var quandlResponse *QuandlResponse

	err = json.Unmarshal(body, &quandlResponse)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//fmt.Printf("%s\n", quandlResponse)
	return quandlResponse, nil
}

// GetData gets Quandl data for a particular identifier and a date range.
// You can optionally set the auth token before running this function so that you
// can make unlimited API calls instead of being limited to 500/day.
func GetData(identifier string, startDate string, endDate string) (*QuandlResponse, error) {
	url := assembleURLwithDates(identifier, startDate, endDate)

	return getDataFromURL(url)
}

// GetAllHistory is similar to GetData except that it does not restrict a date range
func GetAllHistory(identifier string) (*QuandlResponse, error) {
	url := assembleURLwithoutDates(identifier)

	return getDataFromURL(url)
}

// GetTimeSeriesColumn returns the data from the Quandl response for a particular column.
// For some series, particularly stock data, multiple columns are returned. Using this
// method you can specify the specific column to extract.
func (q *QuandlResponse) GetTimeSeriesColumn(column string) []float64 {
	_, data := q.GetTimeSeries(column)

	return data
}

// GetTimeSeriesData returns the most relevant data column from the Quandl response.
// In many cases you will not necessarily know beforehand what type of data is being
// requested and therefore cannot determine if it's stock data vs. economic data. In
// such cases, you can use this function to grab the column that is most likely relevant.
// The method also returns the most relevant
func (q *QuandlResponse) GetTimeSeriesData() ([]float64, string) {
	column := q.getLikelyDataColumnName()

	_, data := q.GetTimeSeries(column)

	return data, column
}

// GetTimeSeriesDate returns the series of dates in the time series
func (q *QuandlResponse) GetTimeSeriesDate() []string {
	column := q.getLikelyDataColumnName()

	date, _ := q.GetTimeSeries(column)

	return date
}

// GetTimeSeries returns a date vector and the value vector for a particular
// column in the QuandlResponse
func (q *QuandlResponse) GetTimeSeries(column string) ([]string, []float64) {
	if q.Data == nil {
		return nil, nil
	}

	dataArray := q.Data.([]interface{})

	dateVector := make([]string, 0, len(dataArray))
	dataVector := make([]float64, 0, len(dataArray))
	dateColumnNum := q.getColumnNum("Date")
	dataColumnNum := q.getColumnNum(column)

	for k, v := range dataArray {
		switch vv := v.(type) {
		case []interface{}:
			// Check that 0 is the date
			switch vv[dateColumnNum].(type) {
			case string:
				dateVector = append(dateVector, vv[dateColumnNum].(string))
			default:
				fmt.Printf("Problem reading %q as a string.\n", vv[0])
				return nil, nil
			}

			// Match the right column with the requested column
			switch vv[dataColumnNum].(type) {
			case float64:
				dataVector = append(dataVector, vv[dataColumnNum].(float64))
			default:
				fmt.Printf("Problem reading %q as a float64.\n", vv[0])
				return nil, nil
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
			return nil, nil
		}
	}

	return dateVector, dataVector
}

// getLikelyDataColumnName finds the column most likely to be the "data"
// column. It either uses adjusted close or just takes the last column in the series
func (q *QuandlResponse) getLikelyDataColumnName() string {
	adjustedCloseColumn := q.getColumnNum("Adj. Close")

	if len(q.Columns) < 1 {
		return "N/A"
	} else if adjustedCloseColumn == -1 {
		return q.Columns[len(q.Columns)-1]
	} else {
		return q.Columns[adjustedCloseColumn]
	}
}

// getColumnNum returns the column number associated with a particular column name.
// It returns -1 if the column is not found.
func (q *QuandlResponse) getColumnNum(column string) int {
	for i, v := range q.Columns {
		if v == column {
			return i
		}
	}

	return -1
}

// Search executes a query against the Quandl API and returns the JSON object
// as a byte stream. In future releases of this Go (golang) Quandl package
// this will return a native object instead of the json
func Search(query string) ([]byte, error) {
	query = strings.Replace(query, " ", "+", -1)

	url := assembleQueryURL(query)

	body, err := readBytesFromUrl(url)

	return body, err
}

func loadPipeDelimited(url string) [][]string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = '|'

	records, _ := reader.ReadAll()

	return records
}

type linefeedConverter struct {
	r *bufio.Reader
}

func newLinefeedConverter(r io.Reader) io.Reader {
	return linefeedConverter{bufio.NewReader(r)}
}

func (r linefeedConverter) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	if err != nil {
		return n, err
	}
	b = b[:n]
	for i := range b {
		if b[i] == '\r' {
			var next byte
			if j := i + 1; j < len(b) {
				next = b[j]
			} else {
				next, err = r.r.ReadByte()
				if err == nil {
					r.r.UnreadByte()
				}
			}
			if next != '\n' {
				b[i] = '\n'
			}
		}
	}
	return n, err
}

func loadCSVMac(url string) [][]string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer resp.Body.Close()

	reader := csv.NewReader(newLinefeedConverter(resp.Body))

	records, _ := reader.ReadAll()

	return records
}

func loadCSV(url string) [][]string {
	//file, err := os.Open(fileName)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	// for i := 0; true; i++ {
	// 	record, err := reader.Read()

	// 	if err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		log.Fatal(err)
	// 		return nil
	// 	} else {
	// 		if i <= 1 {
	// 			fmt.Printf("%q\n", record)
	// 		}
	// 	}
	// }

	records, _ := reader.ReadAll()

	return records
}

func getCsvColumnNumber(csvArray [][]string, column string) int {
	for i, v := range csvArray[0][:] {
		if v == column {
			return i
		}
	}

	return -1
}

// Get various lists
func extractColumns(csvArray [][]string, column1 int, column2 int, skipHeader bool) ([]string, []string) {
	identifier := make([]string, 0, len(csvArray))
	description := make([]string, 0, len(csvArray))

	for i, v := range csvArray[:] {
		for j, v2 := range v {
			// If the header is included then we use the first row
			// otherwise we skip the first row
			if (!skipHeader && i == 0) || i > 0 {
				if j == column1 {
					identifier = append(identifier, v2)
				} else if j == column2 {
					description = append(description, v2)
				}
			}
		}
	}

	return identifier, description
}

// GetAllSecurityList gets all the security identifiers and descriptions
func GetAllSecurityList() ([]string, []string) {

	identifier, description := GetStockList()

	tempIdentifier, tempDescription := GetStockTickerList()
	identifier = append(identifier, tempIdentifier...)
	description = append(description, tempDescription...)

	tempIdentifier, tempDescription = GetETFList()
	identifier = append(identifier, tempIdentifier...)
	description = append(description, tempDescription...)

	tempIdentifier, tempDescription = GetETFTickerList()
	identifier = append(identifier, tempIdentifier...)
	description = append(description, tempDescription...)

	tempIdentifier, tempDescription = GetStockIndexList()
	identifier = append(identifier, tempIdentifier...)
	description = append(description, tempDescription...)

	tempIdentifier, tempDescription = GetCommoditiesList()
	identifier = append(identifier, tempIdentifier...)
	description = append(description, tempDescription...)

	tempIdentifier, tempDescription = GetBitcoinList()
	identifier = append(identifier, tempIdentifier...)
	description = append(description, tempDescription...)

	return identifier, description
}

// GetStockList gets all the Quandl stock codes and descriptions
func GetStockList() ([]string, []string) {
	list := loadCSV(quandlStockWikiList)

	return extractColumns(list, 0, 1, true)
}

// GetStockTickerList gets all the Quandl codes and tickers
func GetStockTickerList() ([]string, []string) {
	list := loadCSV(quandlStockList)

	return extractColumns(list, 2, 0, true)
}

// GetETFList gets all the Quandl codes and ETF descriptions
func GetETFList() ([]string, []string) {
	list := loadCSV(etfList)

	return extractColumns(list, 1, 2, true)
}

// GetETFTickerList gets all the Quandl codes and ETF tickers
func GetETFTickerList() ([]string, []string) {
	list := loadCSV(etfList)

	return extractColumns(list, 1, 0, true)
}

// GetStockIndexList gets all the Quandl codes and stock index descriptions
func GetStockIndexList() ([]string, []string) {
	list := loadCSV(stockIndexList)

	return extractColumns(list, 1, 2, true)
}

// GetCommoditiesList gets all the Quandl codes and commodity descriptions
func GetCommoditiesList() ([]string, []string) {
	list := loadCSV(commoditiesList)

	return extractColumns(list, 1, 0, true)
}

// GetBitcoinList just returns http://www.quandl.com/api/v1/datasets/BITCOIN/BITSTAMPUSD
func GetBitcoinList() ([]string, []string) {
	identifier := make([]string, 1, 1)
	description := make([]string, 1, 1)

	identifier[0] = "BITCOIN/BITSTAMPUSD"
	description[0] = "Bitcoin Exchange Rate (BTC vs. USD) on Bitstamp"
	return identifier, description
}

// Data item identifier
// GetAllDataItemList

// GetFinancialRatiosList returns the list of Damordoran financial ratios. Currently
// this list is hard-coded into the file because Quandl does not provide a file from
// where to read. A caveat about using these is that you need to append the ticker
// in a particular way to use these ratios
func GetFinancialRatiosList() ([]string, []string) {
	test := [][]string{{"FLOAT", "Number of Shares Outstanding"},
		{"INSIDER", "Insider Holdings"},
		{"CAPEX", "Capital Expenditures"},
		{"NET_MARG", "Net Margin"},
		{"INV_CAP", "Invested Capital"},
		{"P_S", "Price to Sales Ratio"},
		{"ROC", "Return on Capital"},
		{"STOCK_PX", "Stock Price"},
		{"MKT_DE", "Market Debt to Equity Ratio"},
		{"CORREL", "Correlation with the Market"},
		{"PE_FWD", "Forward PE Ratio"},
		{"REV_GRO", "Previous Year Growth in Revenues"},
		{"EBIT_1T", "EBIT for Previous Period"},
		{"DIV", "Dividends"},
		{"EPS_FWD", "Forward Earnings Per Share"},
		{"CHG_NCWC", "Change in Non-Cash Working Capital"},
		{"CASH_FV", "Cash as Percentage of Firm Value"},
		{"INST_HOLD", "Institutional Holdings"},
		{"EFF_TAX", "Effective Tax Rate"},
		{"CASH_ASSETS", "Cash as Percentage of Total Assets"},
		{"FIXED_TOT", "Ratio of Fixed Assets to Total Assets"},
		{"BETA_VL", "Value Line Beta"},
		{"BV_ASSETS", "Book Value of Assets"},
		{"BV_EQTY", "Book Value of Equity"},
		{"FCFF", "Free Cash Flow to Firm"},
		{"CASH_REV", "Cash as Percentage of Revenues"},
		{"MKT_CAP", "Market Capitalization"},
		{"EFF_TAX_INC", "Effective Tax Rate on Income"},
		{"EV_SALES", "EV To Sales Ratio"},
		{"TOT_DEBT", "Total Debt"},
		{"INTANG_TOT", "Ratio of Intangible Assets to Total Assets"},
		{"PE_G", "PE to Growth Ratio"},
		{"REINV_RATE", "Reinvestment Rate"},
		{"BOOK_DC", "Book Debt to Capital Ratio"},
		{"EPS_GRO_EXP", "Expected Growth in Earnings Per Share"},
		{"EV_EBIT", "EV to EBIT Ratio"},
		{"PE_CURR", "Current PE Ratio"},
		{"MKT_DC", "Market Debt to Capital Ratio"},
		{"NCWC_REV", "Non-Cash Working Capital as Percentage of Revenues"},
		{"REV_12M", "Trailing 12-month Revenues"},
		{"REV_GRO_EXP", "Expected Growth in Revenues"},
		{"REV_TRAIL", "Trailing Revenues"},
		{"ROE", "Return on Equity"},
		{"EV_EBITDA", "EV to EBITDA Ratio"},
		{"EBITDA", "Earnings Before Interest Taxes Depreciation and Amortization"},
		{"BETA", "3-Year Regression Beta"},
		{"DEPREC", "Depreciation"},
		{"EV_SALESTR", "EV to Trailing Sales Ratio"},
		{"EPS_GRO", "Growth in Earnings Per Share"},
		{"P_BV", "Price to Book Value Ratio"},
		{"NET_INC_TRAIL", "Trailing Net Income"},
		{"PE_TRAIL", "Trailing PE Ratio"},
		{"OP_MARG", "Pre-Tax Operating Margin"},
		{"FIRM_VAL", "Firm Value"},
		{"STDEV", "3-year Standard Deviation of Stock Price"},
		{"TRAD_VOL", "Trading Volume"},
		{"CASH", "Cash"},
		{"DIV_YLD", "Dividend Yield"},
		{"REV_LAST", "Revenues"},
		{"NET_INC", "Net Income"},
		{"EV_BV", "EV to Book Value Ratio"},
		{"REINV", "Reinvestment Amount"},
		{"EBIT", "Earnings Before Interest and Taxes"},
		{"EV_CAP", "EV to Invested Capital Ratio"},
		{"PAYOUT", "Payout Ratio"},
		{"HILO", "Hi-Lo Risk"},
		{"ALLFINANCIALRATIOS", "All Financial Ratios"},
		{"SGA", "Sales General and Administration Expenses"},
		{"EV", "Enterprise Value"},
		{"NCWC", "Non-Cash Working Capital"},
	}

	return extractColumns(test, 0, 1, false)
}

// Economic data (doesn't pertain to a particular security)
// GetEconomicDataList
func GetEconomicDataList() ([]string, []string) {
	list := loadPipeDelimited(economicData)

	identifier, description := extractColumns(list, 0, 1, true)

	return prependList("FRED/", identifier), description
}

// Index membership
// GetSP500Constituents
func GetSP500Constituents() ([]string, []string) {
	//fmt.Printf("%s\n", spxConstituents)
	list := loadCSVMac(spxConstituents)
	/*need to prepend col 0 ticker with WIKI*/

	identifier, description := extractColumns(list, 0, 2, true)

	return prependList("WIKI/", identifier), description
}

// GetDowConstituents
func GetDowConstituents() ([]string, []string) {
	list := loadCSV(dowConstituents)
	/*need to prepend col 0 ticker with WIKI*/

	identifier, description := extractColumns(list, 0, 2, true)

	return prependList("WIKI/", identifier), description
}

// GetNasdaqCompositeConstituents
func GetNasdaqCompositeConstituents() ([]string, []string) {
	list := loadCSV(nasdaqCompositeConstituents)
	/*need to prepend col 0 ticker with WIKI*/

	identifier, description := extractColumns(list, 0, 2, true)

	return prependList("WIKI/", identifier), description
}

// GetNasdaq100Constituents
func GetNasdaq100Constituents() ([]string, []string) {
	list := loadCSV(nasdaq100Constituents)
	/*need to prepend col 0 ticker with WIKI*/

	identifier, description := extractColumns(list, 0, 2, true)

	return prependList("WIKI/", identifier), description
}

// GetFTSE100Constituents
func GetFTSE100Constituents() ([]string, []string) {
	list := loadCSV(ftse100Constituents)

	identifier, description := extractColumns(list, 1, 2, true)

	return identifier, description
}

// Sector mappings
// GetSP500SectorMappings
func GetSP500SectorMappings() ([]string, []string) {
	list := loadCSVMac(spxConstituents)

	identifier, description := extractColumns(list, 0, 3, true)

	return prependList("WIKI/", identifier), description
}

func prependList(prefix string, list []string) []string {
	newList := make([]string, len(list), len(list))

	for i, v := range list {
		newList[i] = fmt.Sprintf("%s%s", prefix, v)
	}

	return newList
}
