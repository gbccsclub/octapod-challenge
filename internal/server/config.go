package server

import "sync"

type Config struct {
	mu sync.Mutex

	// Loop
	updateInterval  int
	timeoutInterval int

	// Maze
	mazeSize int
}

func NewConfig() *Config {
	return &Config{
		updateInterval:  30,
		timeoutInterval: 10,
		mazeSize:        10,
	}
}

func (c *Config) Get() (int, int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.updateInterval, c.timeoutInterval, c.mazeSize
}

func (c *Config) Set(updateInterval, timeoutInterval, mazeSize int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.updateInterval = updateInterval
	c.timeoutInterval = timeoutInterval
	c.mazeSize = mazeSize
}
