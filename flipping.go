package main
import (
  "fmt"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "sort"
  "os"
  "strconv"
)
type Item struct {
  name string
  margin int
  buyAvg int
  buyQty int
  sellAvg int
  sellQty int
}
func main() {
  url := fmt.Sprintf("https://rsbuddy.com/exchange/summary.json")
  margin, err := strconv.Atoi(os.Args[1])
  minPrice, err := strconv.Atoi(os.Args[2])
  maxPrice, err := strconv.Atoi(os.Args[3])
  response, err := http.Get(url)
  if (err != nil) {
    fmt.Printf("bad request idiot %s\n", err)
  }
  //data is json blob
  data, err := ioutil.ReadAll(response.Body)
  if (err != nil) {
    fmt.Printf("bad request idiot %s\n", err)
  }
  //map json to Item
  var responseObj map[string]interface{}
  json.Unmarshal([]byte(data), &responseObj)
  items := []Item{}
  for _, record := range responseObj {
    if rec, ok := record.(map[string] interface{}); ok {
      var name string
      var buyAvg int
      var sellAvg int
      var buyQty int
      var sellQty int
      for key, value := range rec {
        switch key {
        case "name":
          name = value.(string)
        case "buy_average":
          buyAvg = int(value.(float64))
        case "sell_average":
          sellAvg = int(value.(float64))
        case "buy_quantity":
          buyQty = int(value.(float64))
        case "sell_quantity":
          sellQty = int(value.(float64))
        default:
        }
      }
      if ((buyAvg > minPrice && buyAvg < maxPrice) && (sellAvg - buyAvg > margin)) {
        var newItem Item
        newItem.name = name
        newItem.margin = (sellAvg - buyAvg)
        newItem.sellAvg = sellAvg
        newItem.buyAvg = buyAvg
        newItem.buyQty = buyQty
        newItem.sellQty = sellQty
        items = append(items, newItem)
      }
    }
  }
  sort.Slice(items[:], func(i, j int) bool { return items[i].margin < items[j].margin })

  asciiArt := 
  `
                 ______  _  _                _               
                |  ____|| |(_)              (_)              
                | |__   | | _  _ __   _ __   _  _ __    __ _ 
                |  __|  | || || '_ \ | '_ \ | || '_ \  / _  |
                | |     | || || |_) || |_) || || | | || (_| |
                |_|     |_||_|| .__/ | .__/ |_||_| |_| \__, |
                              | |    | |                __/ |
                              |_|    |_|               |___/
 `

  fmt.Println(asciiArt)
  fmt.Println("\n")
  fmt.Printf("%30s %12s %12s %10s %12s %10s", "Name", "Margin", "Buy Avg", "Buy Qty", "Sell Avg", "Sell Qty\n")
  for i:= 0; i < len(items); i++ {
    fmt.Printf("%30s %12d %12d %10d %12d %10d",
      items[i].name,
      items[i].margin,
      items[i].buyAvg,
      items[i].buyQty,
      items[i].sellAvg,
      items[i].sellQty,)
    fmt.Println("")
  }
}