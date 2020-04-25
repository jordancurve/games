// Calculate the number of legal 60-card decks in various Magic the Gathering formats.
//
// Results with 2020-04-25 mtgjson data:
// standard: 2.8e+115 (28017195940711795642224306590465619699273632089594591292106829401511308382259461554677256729310446278759469882348260)
//   modern: 4.98e+166 (49799283634205261025534897730492900253370806192021078346428925761716941521872746267731856884737044867103520340034751129844703497378665734089208354517494440005273285600)
//   legacy: 1.08e+176 (108119704249563116925939634099939024539955247156059140319826487063625851813543756680087658510979368998064312008713451041188333984045271321135762062624855683134524655130937344360)
//  vintage: 1.27e+176 (126753356042903584001232906419110909738178423684812674067967602782643359734092761257198197087118519533610041355670823852576779686214678070342273185637378196874442478343283360005)
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
				if strings.HasPrefix(c.Type, "Basic Land") {
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
