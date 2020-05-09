// Calculate the number of legal 60-card decks in various Magic the Gathering formats.
//
// Using card info from mtgjson version 4.6.3+20200508 (21045 cards):
// standard: 2.8e+115 (28017195940711795642224306590465619699273632089594591292106829401511308382259461554677256729310446278759469882348260)
//  modern: 4.98e+166 (49799283634328953488916231300951000515942374393940843701494598357112083133553976290335263999751254915957305065677337567872398205625101498459743279214358537667604262540)
//  legacy: 1.34e+176 (133852855223192304891024530141456489373824712510430387523685941930827100356727769009091108334853139788710933440308169438342872609915064103809573318409385591561945217990022766745)
// vintage: 1.57e+176 (156833188599404102378934775557412079300694483047821805614426250349773224588612121528646566293567650145703494433337572272901095738790598694925969954243159683576785212428005674360)
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

type Card struct {
	Name       string
	Type       string
	Legalities map[string]string
}

func main() {
	if len(os.Args) != 2 || (len(os.Args) >= 2 && os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Fprintf(os.Stderr, "usage: %s path/to/AllCards.json  # from https://mtgjson.com/json/AllCards.json\n", os.Args[0])
		os.Exit(1)
	}
	allCardsPath := os.Args[1]
	mtgJSON, err := ioutil.ReadFile(allCardsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	limits := FormatLimits(mtgJSON)
	for _, f := range []string{"standard", "modern", "legacy", "vintage"} {
		c := LimitedMultiChoose(60, limits[f])
		fmt.Printf("%8s: %.3g (%v)\n", f, new(big.Float).SetInt(c), c)
	}
}

func FormatLimits(mtgJSON []byte) map[string][]int {
	var cards map[string]Card
	if err := json.Unmarshal(mtgJSON, &cards); err != nil {
		panic(err)
	}
	limits := map[string][]int{}
	for _, c := range cards {
		for f, leg := range c.Legalities {
			if _, ok := limits[f]; !ok {
				limits[f] = []int{}
			}
			lim := 0
			if leg == "Legal" {
				if strings.HasPrefix(c.Type, "Basic Land") ||
					c.Name == "Relentless Rats" || c.Name == "Shadowborn Apostle" {
					lim = 1000
				} else {
					lim = 4
				}
			} else if leg == "Restricted" {
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
