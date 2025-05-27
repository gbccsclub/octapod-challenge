package handler

import (
	"gbccsclub/octopod-challenge/internal/model"
	"gbccsclub/octopod-challenge/internal/payload"
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

func (h *OctapodHandler) UpdateAll(tickId string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, octapod := range h.octapods {
		octapod.TryUpdate(tickId)
	}
}

func (h *OctapodHandler) PingAll(tickId string, status payload.Status) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, octapod := range h.octapods {
		position := octapod.GetPosition()
		sensor := pkg.NewSensor(false, false, false, false) // TODO get sensor from maze
		pingMsg := payload.NewPingMessage(tickId, sensor, position, status)
		err := octapod.Ping(pingMsg)
		if err != nil {
			log.Println("Error pinging", octapod.GetId(), err)
			octapod.Disconnect()
			delete(h.octapods, octapod.GetId())
			continue
		}
	}
}

func (h *OctapodHandler) HandleJoin(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

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
	if _, ok := h.octapods[id]; ok {
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
		h.mu.Lock()
		defer h.mu.Unlock()
		log.Printf("Octapod %s disconnected, removing from map\n", octapodId)
		delete(h.octapods, octapodId)
	}
	octapod := model.NewOctapod(id, pkg.ZeroVec2(), conn, onDisconnect)

	h.octapods[id] = octapod
	octapod.Run()
}
