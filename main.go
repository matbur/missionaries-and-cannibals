package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/matbur/missionaries-and-cannibals/errors"
)

const (
	MAX = 3
	MIN = 0
)

func main() {
	m := []Move{
		{2, 0},
		{1, 0},
		{1, 1},
		{0, 1},
		{0, 2},
	}

	s := NewState()
	p := NewPath(s)
	t := NewTree(p)

	for t.end == 0 {
		p := t.pop(0)
		if p.end > errors.Error_FEW_K {
			t.add(p)
			continue
		}
		for _, v := range m {
			pp := NewPath(p)
			pp.apply(v)
			if pp.end > errors.Error_FEW_K || pp.end == 0 {
				t.add(pp)
			}
		}
	}

	t.printSummary(os.Stdout)
}

type Move struct {
	m, k int
}

type State struct {
	end     errors.Error
	onRight bool
	m, k    int
}

func NewState() State {
	return State{0, true, MAX, MAX}
}

func (s State) apply(m Move) State {
	if s.end != 0 {
		return s
	}

	if s.onRight {
		if s.m-m.m < MIN {
			s.end = errors.Error_FEW_M
			return s
		}
		if s.k-m.k < MIN {
			s.end = errors.Error_FEW_K
			return s
		}
		s.m -= m.m
		s.k -= m.k

	} else {
		if s.m+m.m > MAX {
			s.end = errors.Error_MANY_M
			return s
		}
		if s.k+m.k > MAX {
			s.end = errors.Error_MANY_K
			return s
		}
		s.m += m.m
		s.k += m.k
	}

	if s.m > MIN && s.k > s.m {
		s.end = errors.Error_EATEN_RIGHT
	}
	if s.m < MAX && s.m > s.k {
		s.end = errors.Error_EATEN_LEFT
	}
	if s.m == MIN && s.k == MIN {
		s.end = errors.Error_FINISHED
	}

	s.onRight = !s.onRight
	return s
}

type Path struct {
	end errors.Error
	tab []State
}

func NewPath(i interface{}) Path {
	tab := []State{}
	switch i.(type) {
	case State:
		tab = append(tab, i.(State))
	case Path:
		tab = append(tab, i.(Path).tab...)
	}
	return Path{0, tab}
}

func (s State) label() string {
	boat := "left"
	if !s.onRight {
		boat = "right"
	}
	return fmt.Sprintf(
		"Left: %dM %dC  |  Right: %dM %dC  |  Boat: %s",
		s.m, s.k, MAX-s.m, MAX-s.k, boat,
	)
}

func moveLabel(from, to State) string {
	dm := from.m - to.m
	dk := from.k - to.k
	if dm > 0 || dk > 0 {
		return fmt.Sprintf("Move %dM + %dC to the right", dm, dk)
	}
	return fmt.Sprintf("Move %dM + %dC to the left", -dm, -dk)
}

func (p Path) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Outcome: %s (%d steps)\n", p.end, len(p.tab)-1)
	for i, state := range p.tab {
		if i == 0 {
			fmt.Fprintf(&b, "  Start: %s\n", state.label())
			continue
		}
		fmt.Fprintf(&b, "  %2d. %s\n", i, moveLabel(p.tab[i-1], state))
		fmt.Fprintf(&b, "      %s\n", state.label())
	}
	return b.String()
}

func (p Path) isIn(s State) bool {
	for _, v := range p.tab {
		if s == v {
			return true
		}
	}
	return false
}

func (p *Path) add(s State) {
	if p.end != 0 {
		return
	}

	if p.isIn(s) {
		p.end = errors.Error_LOOP
		return
	}

	p.tab = append(p.tab, s)
	p.end = s.end
}

func (p *Path) apply(m Move) {
	s := p.getLast().apply(m)
	p.add(s)
}

func (p *Path) getLast() State {
	return p.tab[len(p.tab)-1]
}

type Tree struct {
	end errors.Error
	tab []Path
}

func NewTree(p Path) Tree {
	return Tree{0, []Path{p}}
}

func (t Tree) isIn(p Path) bool {
out:
	for _, path := range t.tab {
		if len(path.tab) != len(p.tab) {
			continue
		}
		for i := 0; i < len(p.tab); i++ {
			if p.tab[i] != path.tab[i] {
				continue out
			}
		}
		return true
	}
	return false
}

func (t *Tree) add(p Path) {
	if t.isIn(p) {
		return
	}

	t.tab = append(t.tab, p)

	for _, v := range t.tab {
		if v.end == 0 {
			t.end = 0
			return
		}
		if t.end < v.end {
			t.end = v.end
		}
	}
}

func (t *Tree) pop(i int) Path {
	p := t.tab[i]
	t.tab = append(t.tab[:i], t.tab[i+1:]...)
	return p
}

func (t Tree) successPaths() []Path {
	var out []Path
	seen := make(map[string]struct{})
	for _, p := range t.tab {
		if p.end != errors.Error_FINISHED {
			continue
		}
		key := p.moveKey()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, p)
	}
	return out
}

func (p Path) moveKey() string {
	var moves []string
	for i := 1; i < len(p.tab); i++ {
		moves = append(moves, moveLabel(p.tab[i-1], p.tab[i]))
	}
	return strings.Join(moves, "|")
}

func (t Tree) printSummary(w io.Writer) {
	failures := make(map[errors.Error]int)
	for _, p := range t.tab {
		if p.end != errors.Error_FINISHED && p.end != 0 {
			failures[p.end]++
		}
	}

	solutions := t.successPaths()

	fmt.Fprintln(w, "Missionaries and Cannibals")
	fmt.Fprintln(w, strings.Repeat("=", 40))
	fmt.Fprintf(w, "Search finished: %s\n", t.end)
	fmt.Fprintf(w, "Paths explored:  %d\n", len(t.tab))
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
	for i, p := range solutions {
		fmt.Fprintf(w, "\n--- Solution %d (%d moves) ---\n", i+1, len(p.tab)-1)
		fmt.Fprint(w, p)
	}
}
