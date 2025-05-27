package pkg

type Vector struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (v Vector) Add(other Vector) Vector {
	return Vec2(v.X+other.X, v.Y+other.Y)
}

func Vec2(x, y int) Vector {
	return Vector{x, y}
}

func ZeroVec2() Vector {
	return Vector{0, 0}
}

func (v Vector) Copy() Vector {
	return Vector{v.X, v.Y}
}
