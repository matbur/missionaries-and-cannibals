package main

import (
	"fmt"
)

const (
	MAX         = 3
	MIN         = 0
	EATEN_RIGHT = 1
	EATEN_LEFT  = 2
	LOOP        = 3
	MANY_M      = 4
	MANY_K      = 5
	FEW_M       = 6
	FEW_K       = 7
	FINISHED    = 8
)

type Move struct {
	m, k int
}

type State struct {
	end     bool
	reason  int
	onRight bool
	m, k    int
}

func NewState() State {
	return State{false, 0, true, MAX, MAX}
}

func (s State) apply(m Move) State {
	if s.end {
		return s
	}

	if s.onRight {
		if s.m-m.m < MIN {
			s.end = true
			s.reason = FEW_M
			return s
		}

		if s.k-m.k < MIN {
			s.end = true
			s.reason = FEW_K
			return s
		}
		s.m -= m.m
		s.k -= m.k

	} else {
		if s.m+m.m > MAX {
			s.end = true
			s.reason = MANY_M
			return s
		}
		if s.k+m.k > MAX {
			s.end = true
			s.reason = MANY_K
			return s
		}
		s.m += m.m
		s.k += m.k
	}

	if s.m > MIN && s.k > s.m {
		s.end = true
		s.reason = EATEN_RIGHT
		return s
	}

	if s.m < MAX && s.m > s.k {
		s.end = true
		s.reason = EATEN_LEFT
		return s
	}

	if s.m == MIN && s.k == MIN {
		s.end = true
		s.reason = FINISHED
		return s
	}

	s.onRight = !s.onRight
	return s
}

type Path struct {
	end    bool
	reason int
	tab    []State
}

func NewPath(s State) Path {
	return Path{false, 0, []State{s}}
}
func NewPath2(p Path) Path {
	pp := Path{false, 0, []State{}}
	pp.tab = append(pp.tab, p.tab...)
	return pp
}

func (p Path) String() string {
	s := fmt.Sprintf("\te: %t %d\n", p.end, p.reason)
	for i, v := range p.tab {
		s += fmt.Sprintf("\t\t%d: %v\n", i, v)
	}
	return s
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
	if p.end {
		return
	}

	p.end = s.end
	p.reason = s.reason

	if p.isIn(s) {
		p.end = true
		p.reason = LOOP
		return
	}

	p.tab = append(p.tab, s)
}

func (p *Path) getLast() State {
	return p.tab[len(p.tab)-1]
}

type Tree struct {
	end    bool
	reason int
	tab    []Path
}

func NewTree(p Path) Tree {
	return Tree{false, 0, []Path{p}}
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

	t.end = true
	for _, v := range t.tab {
		if !v.end {
			t.end = false
			return
		}
	}
	t.reason = FINISHED
}

func (t *Tree) pop(i int) Path {
	p := t.tab[i]
	t.tab = append(t.tab[:i], t.tab[i+1:]...)
	return p
}

func (t Tree) String() string {
	s := fmt.Sprintf("e: %t %d\n", t.end, t.reason)
	for i, v := range t.tab {
		s += fmt.Sprintf("%d:\n%v\n", i, v)
	}
	return s
}

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

	for !t.end {
		l := len(t.tab)
		for i := 0; i < l; i++ {
			p := t.tab[i]
			t.pop(i)
			for _, v := range m {
				pp := NewPath2(p)
				ss := pp.getLast().apply(v)
				pp.add(ss)
				t.add(pp)
			}
		}
	}

	//fmt.Println(t)
	for i, v := range t.tab {
		s := v.getLast()
		if s.m == MIN && s.k == MIN {
			fmt.Println(i, v)
		}
	}
}
