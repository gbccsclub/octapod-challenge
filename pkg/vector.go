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

func (v Vector) Up() Vector {
	return Vec2(v.X, v.Y-1)
}

func (v Vector) Down() Vector {
	return Vec2(v.X, v.Y+1)
}

func (v Vector) Left() Vector {
	return Vec2(v.X-1, v.Y)
}

func (v Vector) Right() Vector {
	return Vec2(v.X+1, v.Y)
}

func Vec2Up() Vector {
	return Vec2(0, -1)
}

func Vec2Down() Vector {
	return Vec2(0, 1)
}

func Vec2Left() Vector {
	return Vec2(-1, 0)
}

func Vec2Right() Vector {
	return Vec2(1, 0)
}
