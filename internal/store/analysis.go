package store

import "time"

type StockHistory struct {
	Symbol string
	Data   []DailyData
}

type DailyData struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

type PortfolioAnalysis struct {
	PortfolioID int
	Stocks      []StockHistory
}
