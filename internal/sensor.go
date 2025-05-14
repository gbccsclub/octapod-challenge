package internal

type Sensor struct {
	Left  bool `json:"left"`
	Right bool `json:"right"`
	Up    bool `json:"up"`
	Down  bool `json:"down"`
}
