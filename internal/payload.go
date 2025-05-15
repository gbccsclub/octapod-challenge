package internal

import "github.com/quartercastle/vector"

type PingMessage struct {
	Sensor   *Sensor       `json:"sensor"`
	Position vector.Vector `json:"position"`
}

type Move string

const (
	Up    Move = "Up"
	Down  Move = "Down"
	Left  Move = "Left"
	Right Move = "Right"
)

type MoveMessage struct {
	Move Move `json:"move"`
}

func (move Move) ToVector() vector.Vector {
	switch move {
	case Up:
		return vector.Vector{0, -1}
	case Down:
		return vector.Vector{0, 1}
	case Left:
		return vector.Vector{-1, 0}
	case Right:
		return vector.Vector{1, 0}
	default:
		return vector.Vector{0, 0}
	}
}

type AuthMessage struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}
