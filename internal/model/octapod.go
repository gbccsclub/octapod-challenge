package model

import (
	"encoding/json"
	"gbccsclub/octopod-challenge/internal/payload"
	"gbccsclub/octopod-challenge/pkg"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Octapod struct {
	mu           sync.Mutex
	moveReceived bool
	moveMsg      *payload.MoveMessage

	id           string
	position     pkg.Vector
	conn         *websocket.Conn
	onDisconnect func(id string)
}

func NewOctapod(id string, position pkg.Vector, conn *websocket.Conn, onDisconnect func(id string)) *Octapod {
	return &Octapod{
		moveReceived: false,
		id:           id,
		position:     position,
		conn:         conn,
		onDisconnect: onDisconnect,
	}
}

func (o *Octapod) Run() {
	go o.readLoop()
}

func (o *Octapod) Ping(pingMsg *payload.PingMessage) error {
	o.mu.Lock()
	o.moveReceived = false
	o.moveMsg = nil
	o.mu.Unlock()
	return o.conn.WriteJSON(pingMsg)
}

func (o *Octapod) TryUpdate(tickId string) bool {
	o.mu.Lock()
	defer o.mu.Unlock()

	if !o.moveReceived || o.moveMsg == nil || o.moveMsg.TickId != tickId {
		log.Printf("Move not received for %s\n", o.id)
		return false
	}

	log.Printf("Applying move %s to %s\n", o.moveMsg.MoveDirection, o.id)
	o.position = o.position.Add(o.moveMsg.ToVector())
	o.moveMsg = nil
	o.moveReceived = false
	return true
}

func (o *Octapod) Disconnect() {
	o.mu.Lock()
	defer o.mu.Unlock()
	err := o.conn.Close()
	if err != nil {
		log.Printf("Error closing connection for %s: %v\n", o.id, err)
		return
	}
}

func (o *Octapod) readLoop() {
	for {
		_, data, err := o.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			if o.onDisconnect != nil {
				o.Disconnect()
				o.onDisconnect(o.id)
			}
			return
		}

		var move payload.MoveMessage
		if err := json.Unmarshal(data, &move); err == nil {
			o.mu.Lock()
			log.Printf("Move received from %s: %s\n", o.id, move.MoveDirection)
			o.moveReceived = true
			o.moveMsg = &move
			o.mu.Unlock()
		}
	}
}

func (o *Octapod) GetId() string {
	return o.id
}

func (o *Octapod) GetPosition() pkg.Vector {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.position
}
