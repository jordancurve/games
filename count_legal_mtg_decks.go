// https://twitter.com/jordancurve/status/920773458031792129
// Calculate the number of legal decks (60 main + 15 sideboard) in various Magic the Gathering formats.
// Results (as of 2019-06-03 mtgjson data):
// standard: 3.91e+150 (3909351344688742300288758321903692825050622043009977216137262000738777449027147520175576054766033674875320016580982147811068464179492111937793751377061)
// modern: 6.34e+213 (6344345225455417088402763894754534620716533399551621269696356058530014563492339869963815098988214553066834233744826358391372515964969861064410539437388489273465019960903149138417726512099052635157383820862032724128)
// legacy: 2.23e+226 (22281635943281172571514903146388973895280340280956155138304193048568753462007393574022552332359430674242526210910665693813445221921875705710271424761474554279935502532069257393504139321305640079824750579451523206771374066320167)
// vintage: 2.73e+226 (27315262066946841575155019474967407242182011452087269827072454827203566135910475042414757490111155530029534358488946328446681655375861164754329765774523835547710780172696411669142774653780151105977526260566824553192780299058336)
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
