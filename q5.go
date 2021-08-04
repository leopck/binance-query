package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strconv"
  "sort"
  "time"
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

type ByTrade []Ticket

func (p ByTrade) Len() int           { return len(p) }
func (p ByTrade) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByTrade) Less(i, j int) bool { return p[i].Count < p[j].Count }

func main() {
  d := time.NewTicker(10 * time.Second)
  mychannel := make(chan bool)
  counter := 0
  ask := 0.0
  bid := 0.0
  priceCH := 0.0
  prevpricech := []float64{0.0, 0.0, 0.0, 0.0, 0.0}
  prevaskbid := []float64{0.0, 0.0, 0.0, 0.0, 0.0}
  fmt.Println("Starting... Running every 10 secs interval...")
  // go func() {
  //   time.Sleep(600 * time.Second)

  //   // Setting the value of channel
  //   mychannel <- true
  // }()

  url := "https://api.binance.com/api/v3/ticker/24hr"
  method := "GET"

  for {
    // Select statement
    select {

    // Case statement
    case <-mychannel:
        fmt.Println("Completed!")
        return

    // Case to print current time
    case tm := <-d.C:
      fmt.Println("The Current time is: ", tm)
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
        
      var p ByTrade
      json.Unmarshal([]byte(string(body)), &p)

      sort.Sort(sort.Reverse(p))
      for k := 0; k < len(p) ; k++ {
        if p[k].Symbol[len(p[k].Symbol)-4:] == "USDT" {
          // Bid-Ask Price over 24 hours per Symbol: /api/v3/ticker/24hr : "askPrice": "4.00000200", "bidPrice": "4.00000000"
          // Price Change over 24 hours per Symbol: /api/v3/ticker/24hr : "priceChange": "-94.99999800",
          ask, _ = strconv.ParseFloat(p[k].AskPrice, 12)
          bid, _ = strconv.ParseFloat(p[k].BidPrice, 12)
          priceCH, _ = strconv.ParseFloat(p[k].PriceChange, 12)
          fmt.Println("Symbol: ", p[k].Symbol, " Bid-Ask Price: ", (ask - bid), " Price Change: ", priceCH, " Delta Price Change: ", priceCH - prevpricech[counter], " Delta Bid-Ask Price: ", (ask - bid) - prevaskbid[counter])
          prevpricech[counter] = priceCH
          prevaskbid[counter] = (ask - bid)
          counter++
        }
        if counter >= 5 {
          counter = 0
          break
        }
      }
    }
  }
}