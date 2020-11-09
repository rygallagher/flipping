package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "sort"
  "strconv"
)

type Item struct {
  name string
  margin int
  buyAverage int
  buyQuantity int
  sellAverage int
  sellQuantity int
}

func main() {
  url := fmt.Sprintf("https://rsbuddy.com/exchange/summary.json")

  response, err := http.Get(url)
  margin, err := strconv.Atoi(os.Args[1])
  minPrice, err := strconv.Atoi(os.Args[2])
  maxPrice, err := strconv.Atoi(os.Args[3])
  data, err := ioutil.ReadAll(response.Body)

  if err != nil {
    fmt.Printf("Bad request %s\n", err)
  }

  var responseObj map[string]interface{}
  json.Unmarshal(data, &responseObj)
  var items []Item

  for _, record := range responseObj {
    if rec, ok := record.(map[string] interface{}); ok {
      var name string
      var buyAverage int
      var sellAverage int
      var buyQuantity int
      var sellQuantity int

      for key, value := range rec {
        switch key {
        case "name":
          name = value.(string)
        case "buy_average":
          buyAverage = int(value.(float64))
        case "sell_average":
          sellAverage = int(value.(float64))
        case "buy_quantity":
          buyQuantity = int(value.(float64))
        case "sell_quantity":
          sellQuantity = int(value.(float64))
        default:
        }
      }

      if (buyAverage > minPrice && buyAverage < maxPrice) && (sellAverage-buyAverage > margin) {
        var newItem Item
        newItem.name = name
        newItem.margin = sellAverage - buyAverage
        newItem.sellAverage = sellAverage
        newItem.buyAverage = buyAverage
        newItem.buyQuantity = buyQuantity
        newItem.sellQuantity = sellQuantity
        items = append(items, newItem)
      }
    }
  }

  sort.Slice(items[:], func(i, j int) bool { return items[i].margin < items[j].margin })

  fmt.Println()
  fmt.Printf(
    "%30s %15s %15s %15s %15s %15s %0s",
    "Name",
    "Margin",
    "Buy Average",
    "Buy Quantity",
    "Sell Average",
    "Sell Quantity",
    "\n",
  )

  for i:= 0; i < len(items); i++ {
    fmt.Printf(
      "%30s %15d %15d %15d %15d %15d %0s",
      items[i].name,
      items[i].margin,
      items[i].buyAverage,
      items[i].buyQuantity,
      items[i].sellAverage,
      items[i].sellQuantity,
      "\n",
    )
  }
}