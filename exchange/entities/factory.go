package exchangeEntities

type TradeListenerFactory interface {
	GetListener(exchange string) (TradesListener, error)
}
