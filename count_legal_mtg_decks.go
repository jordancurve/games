// https://twitter.com/jordancurve/status/920773458031792129
// Calculate the number of legal decks (60 main + 15 sideboard) in various Magic the Gathering formats.
// Using card info from mtgjson version 4.6.3+20200508 (21045 cards):
// standard: 3.96e+152 (395697481306288315500482412588185550997575949463159607457342791398956402454575201937306830423839076258993204642893660863880836081092733403218174252555980)
//   modern: 5.3e+216 (5303212499185418908525344852772915685049773837590652362658065020537815355432286385564909790828480055390241936299604049862939115284195644153207212076681276012055210262431562606172865537079398006755301921727478706141752)
//   legacy: 3.27e+228 (3269594370689434972246239696397784575816336438051519915677802141512029916510826903371896498585407893573079813262041191251611293391494569723115805976164312782168475127317700176973707069162557373856126234625020790537917613937349470)
//  vintage: 3.99e+228 (3985786980972339470046639745982013864174118310772370078099163128634390654656416392559271326021684859840974539535937081683655324261736494594244708112735908815730951142524698843844153661953325858096236390658578097533754643413523536)
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
		c := CountDecks(60, 15, limits[f])
		fmt.Printf("%8s: %.3g (%v)\n", f, new(big.Float).SetInt(c), c)
	}
}

func FormatLimits(mtgJSON []byte) map[string][]int {
	var cards map[string]Card
	if err := json.Unmarshal(mtgJSON, &cards); err != nil {
		panic(err)
	}
	limits := map[string][]int{}
	fmt.Printf("%d cards\n", len(cards))
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
