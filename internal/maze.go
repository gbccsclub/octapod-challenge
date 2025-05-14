// AI generated. I will not take responsibility for any damage caused by this code.
package internal

import (
	"github.com/quartercastle/vector"
	"math/rand"
)

type Maze struct {
	Width   int
	Height  int
	cells   [][]bool // true: wall, false: path
	visited [][]bool
}

func NewMaze(width, height int) *Maze {
	m := &Maze{
		Width:   width,
		Height:  height,
		cells:   make([][]bool, width),
		visited: make([][]bool, width),
	}
	for i := range m.cells {
		m.cells[i] = make([]bool, height)
		m.visited[i] = make([]bool, height)
	}
	return m
}

// Generate creates a maze with walls (true) and passages (false)
// Start is at (0,0) and end is at (width-1,height-1)
func (m *Maze) Generate() {
	// First, fill the entire maze with walls
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.cells[x][y] = true
			m.visited[x][y] = false
		}
	}

	// Use depth-first search with backtracking to create paths
	// Start from (1,1) in cell coordinates
	m.carvePassages(1, 1)

	// Create entrance (top-left) and exit (bottom-right)
	m.cells[0][0] = false
	m.cells[1][0] = false
	m.cells[m.Width-1][m.Height-1] = false
	m.cells[m.Width-2][m.Height-1] = false
}

// carvePassages uses depth-first search with backtracking to carve passages
func (m *Maze) carvePassages(x, y int) {
	// Mark the current cell as a passage
	m.cells[x][y] = false

	// Define the four possible directions: N, E, S, W
	directions := []struct{ dx, dy int }{
		{0, -2}, // North
		{2, 0},  // East
		{0, 2},  // South
		{-2, 0}, // West
	}

	// Shuffle the directions for randomness
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	// Try each direction
	for _, dir := range directions {
		newX, newY := x+dir.dx, y+dir.dy

		// Check if the new position is within bounds and unvisited (still a wall)
		if newX >= 0 && newX < m.Width && newY >= 0 && newY < m.Height && m.cells[newX][newY] {
			// Carve a passage by removing the wall between current cell and new cell
			m.cells[x+dir.dx/2][y+dir.dy/2] = false

			// Continue DFS from the new cell
			m.carvePassages(newX, newY)
		}
	}
}

// Print returns a string representation of the maze
// "#" represents walls and " " represents paths
func (m *Maze) Print() string {
	var result string
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.cells[x][y] {
				result += "  " // Wall
			} else {
				result += "# " // Path
			}
		}
		result += "\n"
	}
	return result
}

// Now this is my code.
// Which was auto completed by copilot, but I was actively engaging with it.
// And there's no comments. No comments = Human.

func (m *Maze) IsAvailable(point vector.Vector) bool {
	x := int(point.X())
	y := int(point.Y())
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height && !m.cells[x][y]
}

func (m *Maze) GetSensor(point vector.Vector) *Sensor {
	return &Sensor{
		Up:    m.IsAvailable(point.Add(vector.Vector{0, -1})),
		Right: m.IsAvailable(point.Add(vector.Vector{1, 0})),
		Down:  m.IsAvailable(point.Add(vector.Vector{0, 1})),
		Left:  m.IsAvailable(point.Add(vector.Vector{-1, 0})),
	}
}

func (m *Maze) Visit(position vector.Vector) {
	m.visited[int(position.X())][int(position.Y())] = true
}

func (m *Maze) PrintVisited() string {
	var result string
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if !m.cells[x][y] && m.visited[x][y] {
				result += "# " // Path
			} else {
				result += "  " // Wall
			}
		}
		result += "\n"
	}
	return result
}
