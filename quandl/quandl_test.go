package quandl

import "fmt"

func ExampleSetAuthToken() {
	SetAuthToken("")

	x, _ := GetData("WIKI/AAPL", "2013-01-01", "2013-01-05")
	x.ToDate = "2014-05-22"
	fmt.Printf("%s\n", x)

	// Output:
	// &{WIKI Quandl Open Data AAPL daily 1980-12-12 2014-05-22 [Date Open High Low Close Volume Ex-Dividend Split Ratio Adj. Open Adj. High Adj. Low Adj. Close Adj. Volume] [[2013-01-04 %!s(float64=536.97) %!s(float64=538.63) %!s(float64=525.83) %!s(float64=527) %!s(float64=2.12262e+07) %!s(float64=0) %!s(float64=1) %!s(float64=73.613294063499) %!s(float64=73.840863700808) %!s(float64=72.086109870961) %!s(float64=72.24650533822) %!s(float64=1.485834e+08)] [2013-01-03 %!s(float64=547.88) %!s(float64=549.67) %!s(float64=541) %!s(float64=542.1) %!s(float64=1.26059e+07) %!s(float64=0) %!s(float64=1) %!s(float64=75.108947523158) %!s(float64=75.354338879051) %!s(float64=74.165767339615) %!s(float64=74.316566496868) %!s(float64=8.82413e+07)] [2013-01-02 %!s(float64=553.82) %!s(float64=555) %!s(float64=541.63) %!s(float64=549.03) %!s(float64=2.00185e+07) %!s(float64=0) %!s(float64=1) %!s(float64=75.923262972321) %!s(float64=76.08502934101) %!s(float64=74.252134129678) %!s(float64=75.266601187558) %!s(float64=1.401295e+08)]]}
}

func ExampleGetDataStockPrice() {
	x, _ := GetData("WIKI/AAPL", "2013-01-01", "2013-01-05")
	x.ToDate = "2014-05-22"
	fmt.Printf("%s\n", x)

	// Output:
	// &{WIKI Quandl Open Data AAPL daily 1980-12-12 2014-05-22 [Date Open High Low Close Volume Ex-Dividend Split Ratio Adj. Open Adj. High Adj. Low Adj. Close Adj. Volume] [[2013-01-04 %!s(float64=536.97) %!s(float64=538.63) %!s(float64=525.83) %!s(float64=527) %!s(float64=2.12262e+07) %!s(float64=0) %!s(float64=1) %!s(float64=73.613294063499) %!s(float64=73.840863700808) %!s(float64=72.086109870961) %!s(float64=72.24650533822) %!s(float64=1.485834e+08)] [2013-01-03 %!s(float64=547.88) %!s(float64=549.67) %!s(float64=541) %!s(float64=542.1) %!s(float64=1.26059e+07) %!s(float64=0) %!s(float64=1) %!s(float64=75.108947523158) %!s(float64=75.354338879051) %!s(float64=74.165767339615) %!s(float64=74.316566496868) %!s(float64=8.82413e+07)] [2013-01-02 %!s(float64=553.82) %!s(float64=555) %!s(float64=541.63) %!s(float64=549.03) %!s(float64=2.00185e+07) %!s(float64=0) %!s(float64=1) %!s(float64=75.923262972321) %!s(float64=76.08502934101) %!s(float64=74.252134129678) %!s(float64=75.266601187558) %!s(float64=1.401295e+08)]]}
}

func ExampleGetDataStockFundamentals() {
	x, _ := GetData("DMDRN/MSFT_MKT_CAP", "2000-01-01", "2013-01-05")
	fmt.Printf("%s\n", x)

	// Output:
	// &{DMDRN Damodaran Financial Data MSFT_MKT_CAP annual 2000-06-30 2013-06-30 [Date Market Capitalization] [[2012-06-30 %!s(float64=227057.1)] [2011-06-30 %!s(float64=217062.1)] [2010-06-30 %!s(float64=241362.8)] [2009-06-30 %!s(float64=275188)] [2008-06-30 %!s(float64=172089.1)] [2007-06-30 %!s(float64=336499.4)] [2006-06-30 %!s(float64=294403.6)] [2005-06-30 %!s(float64=280921.5)] [2004-06-30 %!s(float64=292811.8)] [2003-06-30 %!s(float64=293355.8)] [2002-06-30 %!s(float64=300819.4)] [2001-06-30 %!s(float64=364524.4)] [2000-06-30 %!s(float64=281947.4)]]}
}

func ExampleGetTimeSeriesDate() {
	q, _ := GetData("DMDRN/MSFT_MKT_CAP", "2000-01-01", "2013-01-05")

	dates := q.GetTimeSeriesDate()

	fmt.Printf("%q\n", dates)

	// Output:
	// ["2012-06-30" "2011-06-30" "2010-06-30" "2009-06-30" "2008-06-30" "2007-06-30" "2006-06-30" "2005-06-30" "2004-06-30" "2003-06-30" "2002-06-30" "2001-06-30" "2000-06-30"]
}

func ExampleGetTimeSeriesData() {
	q, _ := GetData("DMDRN/MSFT_MKT_CAP", "2000-01-01", "2013-01-05")

	dates, column := q.GetTimeSeriesData()

	fmt.Printf("%s: %q\n", column, dates)

	// Output:
	// Market Capitalization: [%!q(float64=227057.1) %!q(float64=217062.1) %!q(float64=241362.8) %!q(float64=275188) %!q(float64=172089.1) %!q(float64=336499.4) %!q(float64=294403.6) %!q(float64=280921.5) %!q(float64=292811.8) %!q(float64=293355.8) %!q(float64=300819.4) %!q(float64=364524.4) %!q(float64=281947.4)]
}

func ExampleGetAllHistory() {
	x, _ := GetAllHistory("DMDRN/MSFT_MKT_CAP")
	fmt.Printf("%s\n", x)

	// Output:
	// &{DMDRN Damodaran Financial Data MSFT_MKT_CAP annual 2000-06-30 2013-06-30 [Date Market Capitalization] [[2013-06-30 %!s(float64=312297.5)] [2012-06-30 %!s(float64=227057.1)] [2011-06-30 %!s(float64=217062.1)] [2010-06-30 %!s(float64=241362.8)] [2009-06-30 %!s(float64=275188)] [2008-06-30 %!s(float64=172089.1)] [2007-06-30 %!s(float64=336499.4)] [2006-06-30 %!s(float64=294403.6)] [2005-06-30 %!s(float64=280921.5)] [2004-06-30 %!s(float64=292811.8)] [2003-06-30 %!s(float64=293355.8)] [2002-06-30 %!s(float64=300819.4)] [2001-06-30 %!s(float64=364524.4)] [2000-06-30 %!s(float64=281947.4)]]}
}

func Example() {
	x, _ := GetData("BOE/XUDLBK73", "2013-01-01", "2013-01-05")
	x.ToDate = "2014-05-22"
	fmt.Printf("%s\n", x)

	// Output:
	// &{BOE Bank of England XUDLBK73 daily 2005-04-01 2014-05-22 [Date Value] [[2013-01-04 %!s(float64=6.2303)] [2013-01-03 %!s(float64=6.2301)] [2013-01-02 %!s(float64=6.2301)]]}
}

func ExampleLoadCSV() {
	output := loadCSV(quandlStockList)

	fmt.Printf("%q\n", output[0:2][:])

	// Output:
	// [["Ticker" "Stock Name" "Price Code" "Ratios Code" "In Market?"] ["A" "Agilent Technologies" "GOOG/NYSE_A" "DMDRN/A_ALLFINANCIALRATIOS" "Active"]]
}

func ExampleGetStockList() {
	identifier, description := GetStockList()

	fmt.Printf("%q : %q\n", identifier[0], description[0])

	// Output:
	// "WIKI/ACT" : "Actavis, Inc."
}

func ExampleGetAllSecurityList() {
	identifier, description := GetAllSecurityList()

	fmt.Printf("len(identifier)=%v\n%q : %q\n", len(identifier), identifier[20000], description[20000])

	// Output:
	// len(identifier)=21382
	// "GOOG/NYSEARCA_TDTT" : "TDTT"
}

// Commented out for now because the originating files use only \r instead of \r\n

func ExampleGetSP500Constituents() {
	identifier, description := GetSP500Constituents()

	fmt.Printf("len(identifer)=%v\n", len(identifier))
	fmt.Printf("len(description)=%v\n", len(description))
	//fmt.Printf("%q : %q\n", identifier[0], description[0])

	// Output:
	// len(identifer)=500
	// len(description)=500
}

func ExampleGetSP500SectorMappings() {
	identifier, description := GetSP500SectorMappings()

	fmt.Printf("len(identifer)=%v\n", len(identifier))
	fmt.Printf("len(description)=%v\n", len(description))
	//fmt.Printf("%q : %q\n", identifier[0], description[0])

	// Output:
	// len(identifer)=500
	// len(description)=500
}

func ExampleGetEconomicDataList() {
	identifier, description := GetEconomicDataList()

	fmt.Printf("%q : %q\n", identifier[0], description[0])

	// Output:
	// "FRED/LAWFIN" : "Finance and Insurance Wages and Salaries in Louisiana"
}

func ExampleGetFinancialRatiosList() {
	identifier, description := GetFinancialRatiosList()

	fmt.Printf("%q : %q\n", identifier[10], description[10])

	// Output:
	// "PE_FWD" : "Forward PE Ratio"
}

/*func ExampleSearch() {
	body, err := Search("Apple Inc Short Interest")

	if err == nil {
		fmt.Printf("%s", body)
	}

	// Output:
	// Testing
}*/
