package internal

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/quartercastle/vector"
	"golang.org/x/crypto/bcrypt"
)

type Octapod struct {
	Id             string
	Conn           *websocket.Conn
	Position       vector.Vector
	InactiveCount  int
	HashedPassword string
	Sensor         chan *Sensor
	Mutex          sync.Mutex
	Maze           *Maze
}

func NewOctapod(id, password string, conn *websocket.Conn, maze *Maze) *Octapod {
	h := hashPassword(password)
	return &Octapod{
		Id:             id,
		HashedPassword: h,
		Conn:           conn,
		Position:       vector.Vector{0, 0},
		Sensor:         make(chan *Sensor),
		Maze:           maze,
	}
}

func (o *Octapod) VerifyPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(o.HashedPassword), []byte(pw)) == nil
}

func (o *Octapod) Disconnect() {
	o.Mutex.Lock()
	defer o.Mutex.Unlock()
	if o.Conn != nil {
		o.Conn.Close()
		o.Conn = nil
		log.Println("Octapod", o.Id, "disconnected")
	}
}

func (o *Octapod) Run() {
	go o.readPump()
	go o.writePump()
}

func (o *Octapod) readPump() {
	for {
		o.Mutex.Lock()
		conn := o.Conn
		o.Mutex.Unlock()
		if conn == nil {
			return
		}

		typ, msg, err := conn.ReadMessage()
		if err != nil {
			o.Disconnect()
			return
		}
		if typ != websocket.TextMessage {
			continue
		}

		var move MoveMessage
		if err := json.Unmarshal(msg, &move); err != nil {
			log.Println("Invalid move from", o.Id, err)
			continue
		}

		log.Println("Move from [", o.Id, "] to", move.Move)

		o.Mutex.Lock()
		newPos := o.Position.Add(move.Move.ToVector())
		if o.Maze.IsAvailable(newPos) {
			o.Position = newPos
		}
		o.InactiveCount = 0
		o.Mutex.Unlock()
	}
}

func (o *Octapod) writePump() {
	for sensor := range o.Sensor {
		o.Mutex.Lock()
		conn := o.Conn
		pos := o.Position
		o.Mutex.Unlock()

		if conn == nil {
			return
		}
		if sensor == nil {
			continue
		}

		msg := PingMessage{Sensor: sensor, Position: pos}
		b, _ := json.Marshal(msg)
		if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Println("Write error for", o.Id, err)
			o.Disconnect()
			return
		}
	}
}

func hashPassword(pw string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hash error:", err)
		return pw
	}
	return string(h)
}
