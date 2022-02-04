package exchangeEntities

import (
	"time"
)

type RecentTrade struct {
	Id        string    // ID транзакции
	Pair      string    // Торговая пара (из списка выше)
	Price     float64   // Цена транзакции
	Amount    float64   // Объем транзакции
	Side      string    // Как биржа засчитала эту сделку (как buy или как sell)
	Timestamp time.Time // Время транзакции
}

func (r RecentTrade) Valid() bool {
	if r.Id == "" || r.Pair == "" || r.Price == 0 || r.Amount == 0 || r.Side == "" || r.Timestamp.IsZero() {
		return false
	}
	return true
}

type TradesListener interface {
	ListenTrades(pairs ...string) (<-chan RecentTrade, error)
}
