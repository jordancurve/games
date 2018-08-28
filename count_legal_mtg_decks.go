// https://twitter.com/jordancurve/status/920773458031792129
// Calculate the number of legal decks (60 main + 15 sideboard) in various Magic the Gathering formats.
// Results (as of 2018-08-28 mtgjson data):
// Standard: 2.73e+138 (2728414044877596543245135316983591502053279898109794401738089964982213666921148782375339277093721577480802676294272091354015941705203133615)
// Modern: 6.74e+211 (67377103387132724494858143354652721230849485161734377588460693836439107560935986370095039639467653304845967826574353049372791186027226467963036443777258838465657383014209553196703299382271602619066154838600051300)
// Legacy: 1.03e+225 (1025377722494955524224595176452883160186422091083351082837178546710814175550440865163126109519949365723590897055996674473111982912224841941423161329011626936609769434795135660464961378031727848261686561800892379187902580522740)
// Vintage: 1.27e+225 (1267745961625140953651578075538659237893514367924440455356666842434519139164475282720932004168886195541716057239512612026391889779879259481780321384274246411291171860018504968555440708830690575906889602882853027067074169979008)
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
		c := CountDecks(60, 15, limits[f])
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
	main, side, numCards int
}

// CountDecks(M, S, L) returns the number of ways to make a deck with M
// cards in the main deck and S cards in the sideboard where there are len(L)
// cards to choose from, and there can be at most L[I] copies of card I in your
// mainboard and sideboard combined (0 < I < len(L)).
// Examples (mainboard/sideboard):
//   CountDecks(3, 0, []int{1,2,3})=6 (abb abc acc bbc bcc ccd)
//   CountDecks(3, 3, []int{1,2,3})=6 (abb/ccc abc/bcc acc/bbc bbc/acc bcc/abc ccc/abb)
//   CountDecks(3, 1, []int{1,2,3})=12 (abb/c abc/b abc/c acc/b acc/c bbc/a bbc/c bcc/a bcc/b bcc/c ccc/a ccc/b)
//   CountDecks(4, 0, []int{1,2,3})=5 (abbc abcc accc bbcc bccc)
//   CountDecks(4, 1, []int{1,2,3})=8 (abbc/c abcc/b abcc/b accc/b bbcc/a bbcc/c bccc/a bccc/b)
//   CountDecks(4, 2, []int{1,2,3})=5 (abbc/cc abcc/bc accc/bb bbcc/ac bccc/ab)
//   CountDecks(60, 15, []int{75})=1 (the "all islands" example)
func CountDecks(numMain, numSide int, limit []int) *big.Int {
	return _countDecks(numMain, numSide, limit, map[key]*big.Int{})
}

func _countDecks(numMain, numSide int, limit []int, cache map[key]*big.Int) *big.Int {
	if numMain+numSide == 0 {
		return big.NewInt(1)
	}
	if len(limit) == 0 {
		return big.NewInt(0)
	}
	key := key{numMain, numSide, len(limit)}
	if val, ok := cache[key]; ok {
		return val
	}
	sum := big.NewInt(0)
	for m := 0; m <= numMain && m <= limit[0]; m++ {
		for s := 0; s <= numSide && m+s <= limit[0]; s++ {
			sum.Add(sum, _countDecks(numMain-m, numSide-s, limit[1:], cache))
		}
	}
	cache[key] = sum
	return sum
}
