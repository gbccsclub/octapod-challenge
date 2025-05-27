package pkg

type Sensor struct {
	Left  bool `json:"left"`
	Right bool `json:"right"`
	Up    bool `json:"up"`
	Down  bool `json:"down"`
}

func NewSensor(left, right, up, down bool) *Sensor {
	return &Sensor{left, right, up, down}
}
