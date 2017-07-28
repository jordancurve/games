// Calculate the number of legal 75-card decks (60 main deck + 15 sideboard) in various Magic the Gathering formats.
// Results (as of 2017-07-28 mtgjson data):
// Standard: 3.01e+152 (301392966750935224627298406818554052849672810773182479866521612789574537385738326677015531414716511718081165380766952246151140442798522064088841448872000)
// Modern: 2.67e+209 (267288690384613372671663162278241780351196965911279109157717273230325621580562675437659710771377059717771631912904871837546012977270113499191074738887699873916616847170601162912271662668197667947811197217074800)
// Legacy: 1.02e+223 (10243361758464765987377094127212880951843990108913523732373068664535987345931468562282781541568129106632398078903666735021036840698409563318420556365566488735448528259819036265230816325253200105484947880576210519165425515520)
// Vintage: 1.27e+223 (12720293787797022158089718710201114942575244255951494583136480109240791401204724368850245018448000631217824432575745007433726128618564056906640516264914838668934899721150373902976392246408021625533684061071939564677484518960)
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
		c := LimitedMultiChoose(75, limits[f])
		b := new(big.Int).Binomial(75, 15)
		c.Mul(c, b)
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
