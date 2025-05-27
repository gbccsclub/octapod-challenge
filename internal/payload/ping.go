package payload

import (
	"gbccsclub/octopod-challenge/pkg"
)

type Status string

const (
	Exploring Status = "Explore"
	Solving   Status = "Solve"
	Solved    Status = "Solved"
)

type PingMessage struct {
	TickId   string      `json:"tickId"`
	Sensor   *pkg.Sensor `json:"sensor"`
	Position pkg.Vector  `json:"position"`
	Status   Status      `json:"status"`
}

func NewPingMessage(tickId string, sensor *pkg.Sensor, position pkg.Vector, status Status) *PingMessage {
	return &PingMessage{
		TickId:   tickId,
		Sensor:   sensor,
		Position: position,
		Status:   status,
	}
}
