package payload

import (
	"gbccsclub/octopod-challenge/pkg"
)

type MoveDirection string

const (
	Up    MoveDirection = "Up"
	Down  MoveDirection = "Down"
	Left  MoveDirection = "Left"
	Right MoveDirection = "Right"
)

type MoveMessage struct {
	TickId        string        `json:"tickId"`
	MoveDirection MoveDirection `json:"moveDirection"`
}

func (m *MoveMessage) IsValid() bool {
	return m.MoveDirection == Up || m.MoveDirection == Down || m.MoveDirection == Left || m.MoveDirection == Right
}

func (m *MoveMessage) ToVector() pkg.Vector {
	switch m.MoveDirection {
	case Up:
		return pkg.Vec2(0, -1)
	case Down:
		return pkg.Vec2(0, 1)
	case Left:
		return pkg.Vec2(-1, 0)
	case Right:
		return pkg.Vec2(1, 0)
	default:
		return pkg.ZeroVec2()
	}
}
