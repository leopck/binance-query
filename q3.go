package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strconv"
  "sort"
	"sync"
  "strings"
)

type Ticket struct {
  Symbol string
  PriceChange string
  PriceChangePercent string
  WeightedAvgPrice string
  PrevClosePrice string
  LastPrice string
  LastQty string
  BidPrice string
  BidQty string
  AskPrice string
  AskQty string
  OpenPrice string
  HighPrice string
  LowPrice string
  Volume string
  QuoteVolume string
  OpenTime int
  CloseTime int
  FirstId int
  LastId int
  Count int
}

type ByVolume []Ticket

type OrderBook struct {
  LastUpdateId int
  Bids [][]string
  Asks [][]string
}

type Results struct {
  Rsym string
  Rmsg string
  Status bool
}

func (p ByVolume) Len() int           { return len(p) }
func (p ByVolume) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByVolume) Less(i, j int) bool { 
  first, _ := strconv.ParseFloat(p[i].Volume, 12)
  second, _ := strconv.ParseFloat(p[j].Volume, 12)
  return first < second
}

func getSymbolsWithURL() (s [5]string) {
  url := "https://api.binance.com/api/v3/ticker/24hr"
  method := "GET"
  bSymbol := [5]string{}

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }
  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
    
  var p ByVolume
  json.Unmarshal([]byte(string(body)), &p)

  sort.Sort(sort.Reverse(p))
  counter := 0
  for k := 0; k < len(p) ; k++ {
    if p[k].Symbol[len(p[k].Symbol)-3:] == "BTC" {
      // Formula Notional Value = Contract units * Spot Price
      // Contract Unit? - 200 bid and ask? or bid quantity? or ask quantity? or combination of bid-ask quantity? or is it trade order?
      // Notional value = Ask Price/Sell Price * Unit of each contract
      fmt.Println("Symbol: ", p[k].Symbol, " Volume: ", p[k].Volume)
      bSymbol[counter] = "https://api.binance.com/api/v3/depth?symbol=" + p[k].Symbol + "&limit=200"
      counter++
    }
    if counter >= 5 {
      counter = 0
      break
    }
  }
  return bSymbol
}

func getSymbolOrderBook200(url string, c chan Results, wg *sync.WaitGroup) {
	defer (*wg).Done()
	res, err := http.Get(url)
  res2 := new(Results)
  res2.Rsym = url
  
	if err != nil {
    res2.Rmsg = "We couldn't reach " + url
    res2.Status = false
		c <- *res2   // pump the result into the channel
	} else {
    body, err2 := ioutil.ReadAll(res.Body)
    if err2 != nil {
      res2.Rmsg = "Failed to parse body"
      res2.Status = false
      c <- *res2
    } else {
      res2.Rmsg = string(body)
      res2.Status = true
      c <- *res2   // pump the result into the channel
    }
	}
}

func getSymbolNotionalValueTop200(urls [5]string) (ru bool) {
  c := make(chan Results)
  var wg sync.WaitGroup
	for _, link := range urls {
		wg.Add(1)   // This tells the waitgroup, that there is now 1 pending operation here
		go getSymbolOrderBook200(link, c, &wg)
	}

	go func() {
		wg.Wait()	// this blocks the goroutine until WaitGroup counter is zero
		close(c)    // Channels need to be closed, otherwise the below loop will go on forever
	}()    // This calls itself

  // this shorthand loop is syntactic sugar for an endless loop that just waits for results to come in through the 'c' channel
	for msg := range c {
    if msg.Status == true {
      var ob OrderBook
      var singleaskvalue []float64
      totalnotionalvalue := 0.0
      counter1 := 0
      ask1 := 0.0
      ask2 := 0.0
      json.Unmarshal([]byte(msg.Rmsg), &ob)
      split := strings.Split(msg.Rsym, "=")
      matched := strings.Split(split[1], "&")
      // Ask[Price, Quantity]
      // Bids [Price, Quantity]
      if len(ob.Asks) > 0 {
        for kl := 0; kl < len(ob.Asks) ; kl++ {
          ask1, _ = strconv.ParseFloat(ob.Asks[kl][0], 12)
          ask2, _ = strconv.ParseFloat(ob.Asks[kl][1], 12)
          singleaskvalue = append(singleaskvalue, ask1 * ask2)
          totalnotionalvalue += singleaskvalue[counter1]
          counter1++
        }
        fmt.Printf("Total Notional Value for %v: %v\n", matched[0], totalnotionalvalue)
      } else {
        fmt.Printf("Order book for %v is empty, please try again\n", matched[0])
      }
   } else {
     fmt.Println(msg.Rmsg)
     return false
   }
	}

  return false
}

func main() {
  // https://api.binance.com/api/v3/depth?symbol=POABTC&limit=200
  // From depth, we can get "bids":[["0.00000069","608780.00000000"] && "asks":[["0.00000070","19434.00000000"]  price, quantity
  // We know that we are getting 200 bids and asks, each order book for each symbol has 200 bids and ask and every bids and ask has quantity price
  // Based on MIN_NOTIONAL = price * quantity: https://github.com/ccxt/ccxt/issues/1972#issuecomment-366834844
  // From Binance Spot: https://www.binance.com/en/trade/BTC_USDT?theme=dark&type=spot
  // We know that we can get contract price and mark price from aggtrade, p = contract price
  sym := getSymbolsWithURL()
  getSymbolNotionalValueTop200(sym)
}