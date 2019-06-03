// Calculate the number of legal 60-card decks in various Magic the Gathering formats.
//
// Results with 2019-06-03 mtgjson data:
// standard: 7.03e+113 (703157299624545053686771888932520980286428279076474012460151327012433705223062879010429007489792352368415680180400)
// modern: 2.29e+164 (229241017427318043159764934629690835542655152031535534015191438353327502711772302958651348359168786657337832400898840998833277507768943622200054517284702035550051800)
// legacy: 2.48e+174 (2476245970516236803157138976198690342673231699980791266540935465097380031382411336013990929809169770898356715035479025800730182913157595543506016140547900979848605735227225255)
// vintage: 2.91e+174 (2914429382569820774637578150121710890809328668604110727990775935135618637873885636206091853018973867369917756238339541481092956757836236427126456480990881788926023871498707668)
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
)

type Card struct {
	Name       string
	Type       string
	Legalities map[string]string
}

func main() {
	limits := FormatLimits("AllCards.json") // from https://mtgjson.com/json/AllCards.json
	for _, f := range []string{"standard", "modern", "legacy", "vintage"} {
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
