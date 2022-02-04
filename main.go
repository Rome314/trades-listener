package main

import (
	"fmt"
	"sync"

	"ws-exchanges-listener/exchange"
	exchangeEntities "ws-exchanges-listener/exchange/entities"
)

func main() {

	// Go Playground doesn't interactive, thus input is here, but it's pretty simple to get it from JSON :)
	input := map[string][]string{
		"poloenix": {"BTC_USDT", "BTC_ETH", "ETH_USDT"},
	}

	factory := exchange.NewTradesListenerFactory()

	wg := &sync.WaitGroup{}
	for exchng, pairs := range input {
		listener, err := factory.GetListener(exchng)
		if err != nil {
			// actually here should not be panic, but for example it,s ok
			panic(err)
		}
		wg.Add(1)
		go func(prefix string, l exchangeEntities.TradesListener, toLister ...string) {
			defer wg.Done()

			trades, e := l.ListenTrades(toLister...)
			if err != nil {
				fmt.Printf("%s | Error: %s\n", prefix, e.Error())
				return
			}
			for t := range trades {
				fmt.Printf("%s | %+v\n", prefix, t)
			}
		}(exchng, listener, pairs...)

	}

	wg.Wait()

}
