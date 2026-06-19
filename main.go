package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/matbur/missionaries-and-cannibals/errors"
)

const bankSize = 3

// All legal boat loads: at most two people, at least one passenger.
var allMoves = []Move{
	{missionaries: 2, cannibals: 0},
	{missionaries: 1, cannibals: 0},
	{missionaries: 1, cannibals: 1},
	{missionaries: 0, cannibals: 1},
	{missionaries: 0, cannibals: 2},
}

func main() {
	tree := newSearchTree(newPath(initialState()))
	tree.search(allMoves)
	tree.printSummary(os.Stdout)
}

type Move struct {
	missionaries int
	cannibals    int
}

// State describes one point in the puzzle.
// leftMissionaries and leftCannibals count people on the left bank.
// boatOnLeft is true when the boat is on the left bank.
type State struct {
	outcome          errors.Error
	boatOnLeft       bool
	leftMissionaries int
	leftCannibals    int
}

func initialState() State {
	return State{
		outcome:          errors.Error_Error_UNKNOWN,
		boatOnLeft:       true,
		leftMissionaries: bankSize,
		leftCannibals:    bankSize,
	}
}

func (s State) rightMissionaries() int { return bankSize - s.leftMissionaries }
func (s State) rightCannibals() int    { return bankSize - s.leftCannibals }

func (s State) isOngoing() bool {
	return s.outcome == errors.Error_Error_UNKNOWN
}

func (s State) apply(move Move) State {
	if !s.isOngoing() {
		return s
	}

	next := s
	if err := next.transfer(move); err != errors.Error_Error_UNKNOWN {
		next.outcome = err
		return next
	}

	next.outcome = next.checkSafety()
	if next.isOngoing() && next.leftMissionaries == 0 && next.leftCannibals == 0 {
		next.outcome = errors.Error_FINISHED
	}

	next.boatOnLeft = !next.boatOnLeft
	return next
}

func (s *State) transfer(move Move) errors.Error {
	if s.boatOnLeft {
		if s.leftMissionaries-move.missionaries < 0 {
			return errors.Error_FEW_M
		}
		if s.leftCannibals-move.cannibals < 0 {
			return errors.Error_FEW_K
		}
		s.leftMissionaries -= move.missionaries
		s.leftCannibals -= move.cannibals
		return errors.Error_Error_UNKNOWN
	}

	if s.leftMissionaries+move.missionaries > bankSize {
		return errors.Error_MANY_M
	}
	if s.leftCannibals+move.cannibals > bankSize {
		return errors.Error_MANY_K
	}
	s.leftMissionaries += move.missionaries
	s.leftCannibals += move.cannibals
	return errors.Error_Error_UNKNOWN
}

func (s State) checkSafety() errors.Error {
	if s.leftMissionaries > 0 && s.leftCannibals > s.leftMissionaries {
		return errors.Error_EATEN_RIGHT
	}
	if s.leftMissionaries < bankSize && s.leftMissionaries > s.leftCannibals {
		return errors.Error_EATEN_LEFT
	}
	return errors.Error_Error_UNKNOWN
}

func (s State) label() string {
	boat := "right"
	if s.boatOnLeft {
		boat = "left"
	}
	return fmt.Sprintf(
		"Left: %dM %dC  |  Right: %dM %dC  |  Boat: %s",
		s.leftMissionaries, s.leftCannibals,
		s.rightMissionaries(), s.rightCannibals(),
		boat,
	)
}

type Path struct {
	outcome errors.Error
	states  []State
}

func newPath(start State) Path {
	return Path{
		outcome: errors.Error_Error_UNKNOWN,
		states:  []State{start},
	}
}

func (p Path) clone() Path {
	states := make([]State, len(p.states))
	copy(states, p.states)
	// Outcome is recomputed as new states are appended (see original NewPath).
	return Path{
		outcome: errors.Error_Error_UNKNOWN,
		states:  states,
	}
}

func (p Path) lastState() State {
	return p.states[len(p.states)-1]
}

func (p Path) contains(state State) bool {
	for _, seen := range p.states {
		if seen == state {
			return true
		}
	}
	return false
}

func (p *Path) append(state State) {
	if !p.isOngoing() {
		return
	}
	if p.contains(state) {
		p.outcome = errors.Error_LOOP
		return
	}
	p.states = append(p.states, state)
	p.outcome = state.outcome
}

func (p Path) isOngoing() bool {
	return p.outcome == errors.Error_Error_UNKNOWN
}

func (p *Path) tryMove(move Move) {
	next := p.lastState().apply(move)
	p.append(next)
}

func moveLabel(from, to State) string {
	dMissionaries := from.leftMissionaries - to.leftMissionaries
	dCannibals := from.leftCannibals - to.leftCannibals
	if dMissionaries > 0 || dCannibals > 0 {
		return fmt.Sprintf("Move %dM + %dC to the right", dMissionaries, dCannibals)
	}
	return fmt.Sprintf("Move %dM + %dC to the left", -dMissionaries, -dCannibals)
}

func (p Path) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Outcome: %s (%d steps)\n", p.outcome, len(p.states)-1)
	for i, state := range p.states {
		if i == 0 {
			fmt.Fprintf(&b, "  Start: %s\n", state.label())
			continue
		}
		fmt.Fprintf(&b, "  %2d. %s\n", i, moveLabel(p.states[i-1], state))
		fmt.Fprintf(&b, "      %s\n", state.label())
	}
	return b.String()
}

func (p Path) moveKey() string {
	moves := make([]string, 0, len(p.states)-1)
	for i := 1; i < len(p.states); i++ {
		moves = append(moves, moveLabel(p.states[i-1], p.states[i]))
	}
	return strings.Join(moves, "|")
}

type SearchTree struct {
	outcome errors.Error
	paths   []Path
}

func newSearchTree(root Path) SearchTree {
	return SearchTree{
		outcome: errors.Error_Error_UNKNOWN,
		paths:   []Path{root},
	}
}

func (t SearchTree) contains(path Path) bool {
	for _, existing := range t.paths {
		if pathsEqual(existing, path) {
			return true
		}
	}
	return false
}

func pathsEqual(a, b Path) bool {
	if len(a.states) != len(b.states) {
		return false
	}
	for i := range a.states {
		if a.states[i] != b.states[i] {
			return false
		}
	}
	return true
}

func (t *SearchTree) add(path Path) {
	if t.contains(path) {
		return
	}
	t.paths = append(t.paths, path)
}

func (t *SearchTree) hasOngoingPaths() bool {
	for _, path := range t.paths {
		if path.isOngoing() {
			return true
		}
	}
	return false
}

func (t *SearchTree) recomputeOutcome() {
	t.outcome = errors.Error_Error_UNKNOWN
	for _, path := range t.paths {
		if path.isOngoing() {
			return
		}
		if t.outcome < path.outcome {
			t.outcome = path.outcome
		}
	}
}

func (t *SearchTree) popFront() Path {
	front := t.paths[0]
	t.paths = t.paths[1:]
	return front
}

// search explores every valid sequence of moves breadth-first.
func (t *SearchTree) search(moves []Move) {
	for t.hasOngoingPaths() {
		current := t.popFront()

		// Terminal paths (success or hard failure) are kept but not expanded.
		if isTerminalOutcome(current.outcome) {
			t.add(current)
			continue
		}

		for _, move := range moves {
			candidate := current.clone()
			candidate.tryMove(move)
			if candidate.isOngoing() || isTerminalOutcome(candidate.outcome) {
				t.add(candidate)
			}
		}
	}
	t.recomputeOutcome()
}

// Invalid move attempts (too few people on a bank) are discarded.
// Everything else is worth keeping in the search tree.
func isTerminalOutcome(outcome errors.Error) bool {
	return outcome > errors.Error_FEW_K
}

func (t SearchTree) successPaths() []Path {
	var solutions []Path
	seen := make(map[string]struct{})
	for _, path := range t.paths {
		if path.outcome != errors.Error_FINISHED {
			continue
		}
		key := path.moveKey()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		solutions = append(solutions, path)
	}
	return solutions
}

func (t SearchTree) printSummary(w io.Writer) {
	failures := make(map[errors.Error]int)
	for _, path := range t.paths {
		if path.outcome != errors.Error_FINISHED && !path.isOngoing() {
			failures[path.outcome]++
		}
	}

	solutions := t.successPaths()

	fmt.Fprintln(w, "Missionaries and Cannibals")
	fmt.Fprintln(w, strings.Repeat("=", 40))
	fmt.Fprintf(w, "Search finished: %s\n", t.outcome)
	fmt.Fprintf(w, "Paths explored:  %d\n", len(t.paths))
	fmt.Fprintf(w, "Solutions found: %d\n", len(solutions))
	if len(failures) > 0 {
		fmt.Fprintln(w, "\nFailed paths:")
		for err, count := range failures {
			fmt.Fprintf(w, "  %-14s %d\n", err, count)
		}
	}

	if len(solutions) == 0 {
		fmt.Fprintln(w, "\nNo solution found.")
		return
	}

	fmt.Fprintln(w, "\nSolutions:")
	for i, path := range solutions {
		fmt.Fprintf(w, "\n--- Solution %d (%d moves) ---\n", i+1, len(path.states)-1)
		fmt.Fprint(w, path)
	}
}
