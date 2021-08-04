package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strconv"
  "sort"
  "time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
  "log"
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

var (
    gt1ba = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_1_Symbol",
        Name:      "bidaskprice",
        Help:      "Bid-Ask Price Changes over 24 hours",
      })
    gt2ba = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_2_Symbol",
        Name:      "bidaskprice",
        Help:      "Bid-Ask Price Changes over 24 hours",
      })
    gt3ba = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_3_Symbol",
        Name:      "bidaskprice",
        Help:      "Bid-Ask Price Changes over 24 hours",
      })
    gt4ba = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_4_Symbol",
        Name:      "bidaskprice",
        Help:      "Bid-Ask Price Changes over 24 hours",
      })
    gt5ba = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_5_Symbol",
        Name:      "bidaskprice",
        Help:      "Bid-Ask Price Changes over 24 hours",
      })

    gt1pc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_1_Symbol",
        Name:      "pricechange",
        Help:      "Price Changes over 24 hours",
      })
    gt2pc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_2_Symbol",
        Name:      "pricechange",
        Help:      "Price Changes over 24 hours",
      })
    gt3pc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_3_Symbol",
        Name:      "pricechange",
        Help:      "Price Changes over 24 hours",
      })
    gt4pc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_4_Symbol",
        Name:      "pricechange",
        Help:      "Price Changes over 24 hours",
      })
    gt5pc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_5_Symbol",
        Name:      "pricechange",
        Help:      "Price Changes over 24 hours",
      })

    gt1deltaaskbid = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_1_Symbol",
        Name:      "delta_askbid",
        Help:      "Delta Price Changes from previous value",
      })
    gt2deltaaskbid = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_2_Symbol",
        Name:      "delta_askbid",
        Help:      "Delta Price Changes from previous value",
      })
    gt3deltaaskbid = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_3_Symbol",
        Name:      "delta_askbid",
        Help:      "Delta Price Changes from previous value",
      })
    gt4deltaaskbid = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_4_Symbol",
        Name:      "delta_askbid",
        Help:      "Delta Price Changes from previous value",
      })
    gt5deltaaskbid = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_5_Symbol",
        Name:      "delta_askbid",
        Help:      "Delta Price Changes from previous value",
      })

    gt1deltapc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_1_Symbol",
        Name:      "delta_pricechange",
        Help:      "Delta Price Change from the previous value",
      })
    gt2deltapc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_2_Symbol",
        Name:      "delta_pricechange",
        Help:      "Delta Price Change from the previous value",
      })
    gt3deltapc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_3_Symbol",
        Name:      "delta_pricechange",
        Help:      "Delta Price Change from the previous value",
      })
    gt4deltapc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_4_Symbol",
        Name:      "delta_pricechange",
        Help:      "Delta Price Change from the previous value",
      })
    gt5deltapc = prometheus.NewGauge(
      prometheus.GaugeOpts{
        Namespace: "Top_5_Symbol",
        Name:      "delta_pricechange",
        Help:      "Delta Price Change from the previous value",
      })
)

func main() {
  d := time.NewTicker(10 * time.Second)
  http.Handle("/metrics", promhttp.Handler())
  gtpc := []func(float64){gt1pc.Set, gt2pc.Set, gt3pc.Set, gt4pc.Set, gt5pc.Set}
  gtba := []func(float64){gt1ba.Set, gt2ba.Set, gt3ba.Set, gt4ba.Set, gt5ba.Set}
  gtdpc := []func(float64){gt1deltapc.Set, gt2deltapc.Set, gt3deltapc.Set, gt4deltapc.Set, gt5deltapc.Set}
  gtdba := []func(float64){gt1deltaaskbid.Set, gt2deltaaskbid.Set, gt3deltaaskbid.Set, gt4deltaaskbid.Set, gt5deltaaskbid.Set}
  
  mychannel := make(chan bool)
  counter := 0
  ask := 0.0
  bid := 0.0
  priceCH := 0.0
  prevpricech := []float64{0.0, 0.0, 0.0, 0.0, 0.0}
  prevaskbid := []float64{0.0, 0.0, 0.0, 0.0, 0.0}
  fmt.Println("Starting... Running every 10 secs interval...")
  // go func() {
  //     time.Sleep(25 * time.Second)
    
  //     // Setting the value of channel
  //     mychannel <- true
  //   }()
    
    url := "https://api.binance.com/api/v3/ticker/24hr"
    method := "GET"
    go func() {
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
              gtba[counter]((ask - bid))
              gtpc[counter](priceCH)
              gtdba[counter]((ask - bid) - prevaskbid[counter])
              gtdpc[counter](priceCH - prevpricech[counter])
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
      }()
      
      prometheus.MustRegister(gt1ba)
      prometheus.MustRegister(gt1pc)
      prometheus.MustRegister(gt2ba)
      prometheus.MustRegister(gt2pc)
      prometheus.MustRegister(gt3ba)
      prometheus.MustRegister(gt3pc)
      prometheus.MustRegister(gt4ba)
      prometheus.MustRegister(gt4pc)
      prometheus.MustRegister(gt5ba)
      prometheus.MustRegister(gt5pc)
      prometheus.MustRegister(gt1deltapc)
      prometheus.MustRegister(gt2deltapc)
      prometheus.MustRegister(gt3deltapc)
      prometheus.MustRegister(gt4deltapc)
      prometheus.MustRegister(gt5deltapc)
      prometheus.MustRegister(gt1deltaaskbid)
      prometheus.MustRegister(gt2deltaaskbid)
      prometheus.MustRegister(gt3deltaaskbid)
      prometheus.MustRegister(gt4deltaaskbid)
      prometheus.MustRegister(gt5deltaaskbid)
      fmt.Println("Publishing to Prometheus at port 8080")
      log.Fatal(http.ListenAndServe(":8080", nil))
    }