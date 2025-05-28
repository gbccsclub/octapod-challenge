package model

import (
	"gbccsclub/octopod-challenge/pkg"
)

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	Exploring Status = "Explore"
	Solving   Status = "Solve"
	Solved    Status = "Solved"
	Ended     Status = "Ended"
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
