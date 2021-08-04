package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "sort"
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

  url := "https://api.binance.com/api/v3/ticker/24hr"
  method := "GET"

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
  counter := 0
  for k := 0; k < len(p) ; k++ {
    if p[k].Symbol[len(p[k].Symbol)-4:] == "USDT" {
      fmt.Println("Symbol: ", p[k].Symbol, " Number of Trades: ", p[k].Count)
      counter++
    }
    if counter >= 5 {
      counter = 0
      break
    }
  }
}