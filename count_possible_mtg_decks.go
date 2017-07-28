// Calculate the number of legal 60-card decks in various Magic the Gathering formats.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
)

type Legality struct {
	Format, Legality string
}

type Card struct {
	Name       string
	Type       string
	Legalities []Legality
}

func main() {
	limits := FormatLimits("AllCards-x.json") // from https://mtgjson.com/json/AllCards-x.json.zip
	formats := []string{"Standard", "Modern", "Legacy", "Vintage"}
	for _, f := range formats {
		c := LimitedMultiChoose(60, limits[f])
		fmt.Printf("%8s: %.3g (%v)\n", f, new(big.Float).SetInt(c), c)
	}
}

func FormatLimits(mtgJsonFile string) map[string][]int {
	mtgJson, err := ioutil.ReadFile(mtgJsonFile)
	if err != nil {
		panic(err)
	}
	var cards map[string]Card
	if err := json.Unmarshal(mtgJson, &cards); err != nil {
		panic(err)
	}
	limits := map[string][]int{}
	for _, c := range cards {
		for _, leg := range c.Legalities {
			f := leg.Format
			if _, ok := limits[f]; !ok {
				limits[f] = []int{}
			}
			lim := 0
			if leg.Legality == "Legal" {
				if strings.HasPrefix(c.Type, "Basic Land") {
					lim = 1000
				} else {
					lim = 4
				}
			} else if leg.Legality == "Restricted" {
				lim = 1
			}
			if lim > 0 {
				limits[f] = append(limits[f], lim)
			}
		}
	}
	return limits
}

// Cache key, used to speed up LimitedMultiChooose.
type key struct {
	numToBuy, numProducts int
}

// LimitedMultiChoose(B, L) returns the number of ways to buy B items from a
// store with len(L) products, where item I has only L[I] in stock (0 < I < N).
// For example, LimitedMultiChoose(5, []int{3,4,6})=17, which is the number of
// ways to choose 5 items to buy from a selection of 3 products, where product
// 0 has 3 in stock, product 1 has 4 in stock, and product 2 has 6 in stock.
func LimitedMultiChoose(numToBuy int, numInStock []int) *big.Int {
	return _limitedMultiChoose(numToBuy, numInStock, map[key]*big.Int{})
}

func _limitedMultiChoose(numToBuy int, numInStock []int, cache map[key]*big.Int) *big.Int {
	if numToBuy == 0 {
		return big.NewInt(1)
	}
	if numToBuy < 0 || len(numInStock) == 0 {
		return big.NewInt(0)
	}
	key := key{numToBuy, len(numInStock)}
	if val, ok := cache[key]; ok {
		return val
	}
	sum := big.NewInt(0)
	for i := 0; i <= numInStock[0]; i++ {
		sum.Add(sum, _limitedMultiChoose(numToBuy-i, numInStock[1:], cache))
	}
	cache[key] = sum
	return sum
}
