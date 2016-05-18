package mdp

import (
	"math"
	"testing"
)

func TestValues(t *testing.T) {
	tolerance := 1e-16
	cases := []struct {
		discount float64
		states   MDP
		want     []Value
	}{
		{
			// You flip a coin until it comes up heads. You get a dollar for every flip.
			discount: 1.0,
			states: MDP{
				State{Nature, 1, []Action{{0, 0.5}, {1, 0.5}}},
				State{Nature, 1, []Action{}},
			},
			want: []Value{2, 0},
		},
		{
			// Fair Duel (from The Population Explosion, by Dick Hess)
			// https://groups.yahoo.com/neo/groups/fallible-ideas/conversations/messages/15171
			discount: 1.0,
			states: MDP{
				State{Nature, 0, []Action{{2, 0.4}, {1, 0.6}}}, // Smith shoots
				State{Nature, 0, []Action{{3, 0.6}, {0, 0.4}}}, // Brown shoots
				State{Nature, 1, []Action{}},                   // Smith hits
				State{Nature, 0, []Action{}},                   // Brown hits
			},
			want: []Value{10.0 / 19, 4.0 / 19, 0, 0},
		},
		{
			// Bus Ticket Roulette, from The Population Explosion by Dick Hess
			// (have $2, want $4 version)
			// https://groups.yahoo.com/neo/groups/fallible-ideas/conversations/messages/15127
			discount: 1.0,
			states: MDP{
				State{Nature, 0, []Action{}},                               // 0: $0 (lose)
				State{Player1, 0, []Action{{5, 0}}},                        // 1: $1
				State{Player1, 0, []Action{{6, 0}, {7, 0}}},                // 2: $2
				State{Player1, 0, []Action{{8, 0}}},                        // 3: $3
				State{Nature, 1, []Action{}},                               // 4: $4+ (win)
				State{Nature, 0, []Action{{2, 18.0 / 38}, {0, 20.0 / 38}}}, // 5: have $1, bet $1 on red
				State{Nature, 0, []Action{{3, 18.0 / 38}, {1, 20.0 / 38}}}, // 6: have $2, bet $1 on red
				State{Nature, 0, []Action{{4, 18.0 / 38}, {0, 20.0 / 38}}}, // 7: have $2, bet $2 on red
				State{Nature, 0, []Action{{4, 18.0 / 38}, {2, 20.0 / 38}}}, // 8: have $3, bet $1 on red
			},
			want: []Value{0, (18.0 * 18.0) / (38 * 38), 18.0 / 38, 18.0/38 + (18.0*20.0)/(38*38)},
		},
	}
	for _, c := range cases {
		got := c.states.Values(c.discount, tolerance)
		for i, r := range c.want {
			if r != 0 && math.Abs(float64(r-got[i])) > tolerance {
				t.Errorf("Values(%+v) V[%d]: got %v; want %v", c.states, i, got[i], r)
			}
		}
	}
}
