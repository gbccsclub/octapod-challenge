package pkg

// Sensor represents the sensors on the octapod.
// True if blocked
type Sensor struct {
	Left  bool `json:"left"`
	Right bool `json:"right"`
	Up    bool `json:"up"`
	Down  bool `json:"down"`
}

func (s Sensor) IsBlocked(direction Vector) bool {
	switch direction {
	case Vec2Up():
		return s.Up
	case Vec2Down():
		return s.Down
	case Vec2Right():
		return s.Right
	case Vec2Left():
		return s.Left
	default:
		return true
	}
}

func NewSensor(left, right, up, down bool) *Sensor {
	return &Sensor{left, right, up, down}
}
