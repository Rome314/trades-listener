package exchange

import (
	"fmt"

	"ws-exchanges-listener/exchange/entities"
	"ws-exchanges-listener/poloenix"
)

type factory struct {
}

func NewTradesListenerFactory() exchangeEntities.TradeListenerFactory {
	return &factory{}
}

func (f *factory) GetListener(exchange string) (exchangeEntities.TradesListener, error) {
	switch exchange {
	case "poloenix":
		return poloenix.NewPoloenix()
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", exchange)
	}
}
