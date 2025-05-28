package server

import (
	"gbccsclub/octopod-challenge/internal/model"
	"gbccsclub/octopod-challenge/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type OctapodHandler struct {
	mu       sync.Mutex
	octapods map[string]*model.Octapod
}

func NewOctapodHandler() *OctapodHandler {
	return &OctapodHandler{
		octapods: make(map[string]*model.Octapod),
	}
}

func (oh *OctapodHandler) UpdateAll(maze *Maze) []*model.Octapod {
	oh.mu.Lock()
	defer oh.mu.Unlock()
	solvedOctapods := make([]*model.Octapod, 0)
	for _, octapod := range oh.octapods {
		newPosition := octapod.TryUpdate()
		if maze.IsSolved(newPosition) {
			solvedOctapods = append(solvedOctapods, octapod)
		}
	}
	return solvedOctapods
}

func (oh *OctapodHandler) PingAll(tickId string, status model.Status, maze *Maze) {
	oh.mu.Lock()
	defer oh.mu.Unlock()
	for _, octapod := range oh.octapods {
		position := octapod.GetPosition()
		sensor := maze.GetSensor(position)
		if !maze.IsSolved(position) {
			status = model.Solved
		}

		err := octapod.Ping(tickId, sensor, status)
		if err != nil {
			log.Println("Error pinging", octapod.GetId(), err)
			octapod.Disconnect()
			delete(oh.octapods, octapod.GetId())
			continue
		}
	}
}

func (oh *OctapodHandler) HandleJoin(c *gin.Context) {
	oh.mu.Lock()
	defer oh.mu.Unlock()

	id := c.Query("id")
	if id == "" {
		c.String(400, "Missing id")
		return
	}

	isValid, msg := pkg.IsValidID(id)
	if !isValid {
		c.String(400, msg)
		return
	}

	log.Println("New connection attempt from", id)

	// Check if octapod already exists
	if _, ok := oh.octapods[id]; ok {
		c.String(400, "Octapod already exists")
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("New connection established")

	onDisconnect := func(octapodId string) {
		oh.mu.Lock()
		defer oh.mu.Unlock()
		log.Printf("Octapod %s disconnected, removing from map\n", octapodId)
		delete(oh.octapods, octapodId)
	}
	octapod := model.NewOctapod(id, pkg.ZeroVec2(), conn, onDisconnect)

	oh.octapods[id] = octapod
	octapod.Run()
}

func (oh *OctapodHandler) GetOctapodPositionSet() map[pkg.Vector]string {
	// TODO make sure this lock doesn't mess everything up
	oh.mu.Lock()
	defer oh.mu.Unlock()
	positions := make(map[pkg.Vector]string)
	for _, octapod := range oh.octapods {
		positions[octapod.GetPosition()] = octapod.GetId()
	}
	return positions
}

func (oh *OctapodHandler) GetOctapodCount() int {
	return len(oh.octapods)
}
