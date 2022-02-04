package poloenix

import (
	exchangeEntities "ws-exchanges-listener/exchange/entities"
)

const poloenixWsUrl = "wss://api2.poloniex.com"
const poloenixPublicUrl = "https://poloniex.com/public"

type subMessage struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
}

type poloenix struct {
	// channel_id to pair map
	chanToPair map[float64]string
	// pair to channel_id  map
	pairToChan map[string]float64

	conn       *wsConn
	tradesChan chan exchangeEntities.RecentTrade
}

func NewPoloenix() (e exchangeEntities.TradesListener, err error) {
	p := &poloenix{
		chanToPair: map[float64]string{},
		pairToChan: map[string]float64{},
		tradesChan: make(chan exchangeEntities.RecentTrade),
	}

	// fill the maps with currencies and their ids
	if err = p.initChannels(); err != nil {
		return
	}

	connection, err := getConnection()
	if err != nil {
		return
	}

	p.conn = connection
	go p.listenMessages()

	return p, nil
}

func (p *poloenix) ListenTrades(pairs ...string) (tradesChan <-chan exchangeEntities.RecentTrade, err error) {
	for _, pair := range pairs {
		err = p.conn.WriteJSON(subMessage{
			Command: "subscribe",
			Channel: reversePair(pair),
		})
		if err != nil {
			return
		}
	}

	return p.tradesChan, nil
}

func (p *poloenix) listenMessages() {

	for {
		message := []interface{}{}
		if err := p.conn.ReadJSON(&message); err != nil {
			continue
		}

		channel, ok := message[0].(float64)
		if !ok {
			continue
		}

		pair, ok := p.chanToPair[channel]
		if !ok {
			continue
		}

		data, ok := message[2].([]interface{})
		if !ok {
			continue
		}

		for _, eventRaw := range data {
			event, okk := eventRaw.([]interface{})
			if !okk {
				continue
			}

			eventType, okk := event[0].(string)
			if !okk {
				continue
			}
			// ignore other event types
			if eventType != "t" {
				continue
			}

			trade := getTradeFromEvent(event)
			trade.Pair = reversePair(pair)

			if !trade.Valid() {
				continue
			}
			p.tradesChan <- trade

		}

	}
}
