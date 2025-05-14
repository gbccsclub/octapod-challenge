package internal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/quartercastle/vector"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sync"
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

func NewOctapod(
	id string,
	password string,
	conn *websocket.Conn,
	maze *Maze,
) *Octapod {
	hashedPassword := hashPassword(password)
	return &Octapod{
		Id:             id,
		HashedPassword: hashedPassword,
		Conn:           conn,
		Position:       vector.Vector{0, 0},
		Sensor:         make(chan *Sensor),
		InactiveCount:  0,
		Maze:           maze,
	}
}

func (o *Octapod) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(o.HashedPassword), []byte(password)) == nil
}

func (o *Octapod) Disconnect() {
	if o.Conn == nil {
		return
	}

	err := o.Conn.Close()
	if err != nil {
		// Check if it's a 'use of closed network connection' error
		// This is not a critical error, just means the connection was already closed
		if err.Error() != "use of closed network connection" {
			log.Println("Error closing connection:", err)
		}
	}

	o.Conn = nil
	log.Println("Octopod " + o.Id + " disconnected")
}

func (o *Octapod) Connect(conn *websocket.Conn) {
	o.Conn = conn
}

func (o *Octapod) Run() {
	o.InactiveCount = 0
	go o.readPump()
	go o.writePump()
}

func (o *Octapod) readPump() {
	for {
		o.Mutex.Lock()
		if o.Conn == nil {
			return
		}
		o.Mutex.Unlock()

		messageType, content, err := o.Conn.ReadMessage()
		if err != nil {
			o.Disconnect()
			return
		}

		if messageType != websocket.TextMessage {
			continue
		}

		o.Mutex.Lock()
		if o.Sensor == nil { // Only accept client move when sensor is not nil
			o.Mutex.Unlock()
			continue
		}

		var moveMsg MoveMessage
		err = json.Unmarshal(content, &moveMsg)
		if err != nil {
			log.Println("Octopod [", o.Id, "] sent invalid move message:", err)
			o.Mutex.Unlock()
			continue
		}

		log.Println("Octopod [", o.Id, "] sent move message:", moveMsg.Move)

		var moveVector = moveMsg.Move.ToVector()
		var newPosition = o.Position.Add(moveVector)

		if o.Maze.IsAvailable(newPosition) {
			o.Position = newPosition
		}

		o.InactiveCount = 0
		o.Sensor = nil
		o.Mutex.Unlock()
	}
}

func (o *Octapod) writePump() {
	for {
		o.Mutex.Lock()
		if o.Conn == nil {
			return
		}
		o.Mutex.Unlock()

		// Use select with timeout to avoid blocking indefinitely
		// TODO: test this out
		var sensor *Sensor
		select {
		case s, ok := <-o.Sensor:
			if !ok {
				// Channel was closed
				log.Println("Sensor channel closed for octapod", o.Id)
				return
			}
			sensor = s
		}

		if sensor == nil {
			continue
		}

		pingMsg := PingMessage{
			Sensor:   sensor,
			Position: o.Position,
		}

		pingMsgJson, _ := json.Marshal(pingMsg)
		err := o.Conn.WriteMessage(websocket.TextMessage, pingMsgJson)
		if err != nil {
			// Only log non-standard close errors
			if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) && err.Error() != "use of closed network connection" {
				log.Println("Error writing ping message:", err, "for octapod", o.Id)
			}
			o.Disconnect()
			return
		}
	}
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return password // Fallback to plain password on error
	}
	return string(hash)
}
