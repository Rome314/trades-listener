package poloenix

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ws-exchanges-listener/exchange/entities"
)

// ["t","23110244",1,"2606.60000000","0.00046548",1643917023,"1643917023489"] --> RecentTrade
func getTradeFromEvent(event []interface{}) exchangeEntities.RecentTrade {
	var (
		side, tradeId string
		epochMs       int64
		price, amount float64
	)

	// every field is checked by type to be sure that data is correct
	if trId, okk := event[1].(string); okk {
		tradeId = trId

	}
	if s, okk := event[2].(float64); okk {
		if s == 1 {
			side = "BUY"
		} else {
			side = "SELL"
		}
	}
	if pr, okk := event[3].(string); okk {
		price, _ = strconv.ParseFloat(pr, 64)
	}
	if sz, okk := event[4].(string); okk {
		amount, _ = strconv.ParseFloat(sz, 64)
	}
	if ep, okk := event[6].(string); okk {
		epochMs, _ = strconv.ParseInt(ep, 10, 64)
	}

	trade := exchangeEntities.RecentTrade{
		Id:        tradeId,
		Price:     price,
		Amount:    amount,
		Side:      side,
		Timestamp: time.UnixMilli(epochMs),
	}
	return trade
}

// BTC_USDT --> USDT_BTC
func reversePair(pair string) string {
	splited := strings.Split(pair, "_")
	if len(splited) == 2 {
		return fmt.Sprintf("%s_%s", splited[1], splited[0])
	}
	return pair
}

// This is necessary because ws api provides only channel_id, not human-readable pair
func (p *poloenix) initChannels() error {
	url := fmt.Sprintf("%s?command=returnTicker", poloenixPublicUrl)
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("getting response: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", response.Status)
	}

	bts, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	responseBody := map[string]struct {
		Id float64 `json:"id"`
	}{}

	if err = json.Unmarshal(bts, &responseBody); err != nil {
		return fmt.Errorf("unmarshaling response body: %v", err)
	}

	for pair, ticker := range responseBody {
		p.chanToPair[ticker.Id] = pair
		p.pairToChan[pair] = ticker.Id
	}
	return nil
}
