package cards

import "fmt"

type Pos struct {
	X int
	Y int
}

func (p *Pos) Add(x, y int) *Pos {
	return &Pos{
		X: p.X + x,
		Y: p.Y + y,
	}
}

func (p *Pos) Equals(other *Pos) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p *Pos) String() string {
	return fmt.Sprintf("[%d, %d]", p.X, p.Y)
}

type Move struct {
	From Pos
	To   Pos
}

func (m *Move) String() string {
	return fmt.Sprintf("%v -> %v", m.From, m.To)
}

type MoveHistory struct {
	stack []*Move
}

func (h *MoveHistory) Size() int {
	return len(h.stack)
}

func (h *MoveHistory) String() string {
	return fmt.Sprintf("%v", h.stack)
}

func (h *MoveHistory) Push(m *Move) {
	h.stack = append(h.stack, m)
}

func (h *MoveHistory) Pop() *Move {
	n := len(h.stack) - 1
	m := h.stack[n]

	h.stack = h.stack[:n]

	return m
}
