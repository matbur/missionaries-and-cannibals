package main

import (
	"fmt"
)

const (
	MAX = 3
	MIN = 0

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
	end     int
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
			s.end = FEW_M
			return s
		}
		if s.k-m.k < MIN {
			s.end = FEW_K
			return s
		}
		s.m -= m.m
		s.k -= m.k

	} else {
		if s.m+m.m > MAX {
			s.end = MANY_M
			return s
		}
		if s.k+m.k > MAX {
			s.end = MANY_K
			return s
		}
		s.m += m.m
		s.k += m.k
	}

	if s.m > MIN && s.k > s.m {
		s.end = EATEN_RIGHT
		return s
	}
	if s.m < MAX && s.m > s.k {
		s.end = EATEN_LEFT
		return s
	}
	if s.m == MIN && s.k == MIN {
		s.end = FINISHED
	}

	s.onRight = !s.onRight
	return s
}

type Path struct {
	end int
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

func (p Path) String() string {
	s := fmt.Sprintf("\te: %v\n", p.end)
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
	if p.end != 0 {
		//fmt.Println(p.reason)
		return
	}

	if p.isIn(s) {
		p.end = LOOP
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
	end int
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

func (t Tree) String() string {
	s := fmt.Sprintf("e: %v\n", t.end)
	for i, v := range t.tab {
		s += fmt.Sprintf("%d:\n%v\n", i, v)
	}
	return s
}

func isValid(i int) bool {
	return i < LOOP || i == FINISHED
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

	for t.end == 0 {
		l := len(t.tab)
		for i := l - 1; i > -1; i-- {
			p := t.pop(i)
			for _, v := range m {
				pp := NewPath(p)
				pp.apply(v)
				t.add(pp)
			}
		}
	}

	fmt.Println(t)
}
