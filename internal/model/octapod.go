package model

import (
	"encoding/json"
	"gbccsclub/octopod-challenge/pkg"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Octapod struct {
	mu           sync.Mutex
	moveReceived bool
	moveMsg      *MoveMessage
	sensor       *pkg.Sensor
	tickId       string

	id           string
	position     pkg.Vector
	conn         *websocket.Conn
	onDisconnect func(id string)
}

func NewOctapod(id string, position pkg.Vector, conn *websocket.Conn, onDisconnect func(id string)) *Octapod {
	return &Octapod{
		moveReceived: false,
		moveMsg:      nil,
		sensor:       pkg.NewSensor(true, true, true, true),
		id:           id,
		position:     position,
		conn:         conn,
		onDisconnect: onDisconnect,
	}
}

func (o *Octapod) Run() {
	go o.readLoop()
}

func (o *Octapod) Ping(tickId string, sensor *pkg.Sensor, status Status) error {
	o.mu.Lock()
	o.moveReceived = false
	o.moveMsg = nil
	o.tickId = tickId
	o.sensor = sensor
	o.mu.Unlock()

	pingMsg := NewPingMessage(tickId, sensor, o.position, status)
	return o.conn.WriteJSON(pingMsg)
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

func (o *Octapod) TryUpdate() pkg.Vector {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.moveMsg == nil {
		return o.position
	}

	moveDirection := o.moveMsg.ToVector()

	if o.sensor.IsBlocked(moveDirection) {
		log.Printf("Move blocked for %s: %s\n", o.id, o.moveMsg.MoveDirection)
		return o.position
	}

	log.Printf("Applying moveMsg %s to %s\n", o.moveMsg.MoveDirection, o.id)
	o.position = o.position.Add(moveDirection)
	return o.position
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

		if o.moveReceived {
			log.Printf("Move already received for %s\n", o.id)
			continue
		}
		o.moveReceived = true

		var moveMsg MoveMessage
		err = json.Unmarshal(data, &moveMsg)
		if err != nil {
			log.Printf("Error unmarshalling moveMsg for %s: %v\n", o.id, err)
			continue
		}

		if o.tickId != moveMsg.TickId {
			log.Printf("Tick id mismatch for %s: %s != %s\n", o.id, o.tickId, moveMsg.TickId)
			continue
		}

		o.mu.Lock()
		log.Printf("Move received from %s: %s\n", o.id, moveMsg.MoveDirection)
		o.moveMsg = &moveMsg
		o.mu.Unlock()
	}
}

func (o *Octapod) GetId() string {
	return o.id
}

func (o *Octapod) GetPosition() pkg.Vector {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.position.Copy()
}
