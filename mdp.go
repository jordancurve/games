// Package mdp finds the value of each state in a zero- or one-player Markov
// Decision Process under the optimal policy.  This information can be used to
// determine the value of various zero- or one-player games of chance, and to
// determine the optimal policy for one-player games of chance.
package mdp

import (
	"math"
)

/*
	MDP represents a Markov Decision Process. The underlying structure is an
	edge-labeled, node-labeled directed graph.  In what follows, a node is called
	a state and an edge is called an action.  Each state has a player which must
	be either 0 (nature) or 1 (player).  The actions of a nature-state are
	labelled with the float64 probability (must sum to 1) of nature taking that
	edge.  The probabilities of the actions of a player-state are irrelevant.
	The game starts with a token on some state. If it is a nature state, nature
	chooses an action at random according to the distribution of probabilities of
	outgoing edges. The token moves to the resulting state and the player's score
	increases by the resulting state's reward.  If it is a player state, the
	player chooses an action and the token moves to the resulting state. The
	players score increases by the resulting state's reward.  If a state has no
	actions, the game ends. The player's goal is to maximize the final score.
*/
type Player uint

var negInf = math.Inf(-1)

const (
	Nature  = 0
	Player1 = 1
)

type MDP []State

type Action struct {
	NextState uint
	Prob      float64
}

type Value float64
type State struct {
	Player Player
	Reward Value
	Action []Action
}

// Value Iteration
// http://www.cs.berkeley.edu/~pabbeel/cs287-fa12/slides/mdps-exact-methods.pdf
// The returned array of values represents the value of each state, assuming
// the player plays optimally.  To play optimally, the player should, at each
// turn, select the action leading to the highest-valued state.  The discount
// is the amount to discount rewards from future states, and the tolerance is
// the amount two values must be within to be considered equal for the purposes
// of ending the iteration process.
// The algorithm is a kind of expectimax with loops.
func (states MDP) Values(discount, tolerance float64) []Value {
	γ := Value(discount)
	V := [][]Value{
		make([]Value, len(states)),
		make([]Value, len(states)),
	}
	prev, cur := 0, 1
	for {
		diff := false // Are this iteration's values different than the previous?
		for i, state := range states {
			if len(state.Action) == 0 {
				continue
			}
			switch state.Player {
			case Nature:
				ev := Value(0.0) // Expected value of reward for next state.
				for _, action := range state.Action {
					s := action.NextState
					ev += Value(action.Prob) * (states[s].Reward + γ*V[prev][s])
				}
				V[cur][i] = ev
			case Player1:
				max := negInf // Maximum value of reward for next state.
				for _, action := range state.Action {
					s := action.NextState
					max = math.Max(max, float64(states[s].Reward+γ*V[prev][s]))
				}
				V[cur][i] = Value(max)
			}
			if math.Abs(float64(V[prev][i]-V[cur][i])) > tolerance {
				diff = true
			}
		}
		if !diff {
			break
		}
		prev, cur = cur, prev
	}
	return V[0]
}
