package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var MaxInactive = 2

type Lobby struct {
	Maze         *Maze
	Octapods     map[string]*Octapod
	Mutex        sync.Mutex
	timerRunning bool
}

func NewLobby(width, height int) *Lobby {
	maze := NewMaze(width, height)
	maze.Generate()

	lobby := &Lobby{
		Maze:         maze,
		Octapods:     make(map[string]*Octapod),
		timerRunning: false,
	}

	fmt.Println()
	fmt.Println(maze.Print())
	fmt.Println()

	lobby.StartTimer(5*time.Second, 3*time.Second)
	return lobby
}

func (l *Lobby) StartTimer(duration, timeout time.Duration) {
	if l.timerRunning {
		log.Println("Timer already running.")
		return
	}

	l.timerRunning = true

	go func() {
		timer := time.NewTimer(duration)
		defer timer.Stop()
		isWaitingForTimeout := false
		nextDuration := timeout

		for {
			select {
			case <-timer.C:
				if !isWaitingForTimeout {
					l.Update()
					log.Println("Sensor data pinged (Update)")
					nextDuration = timeout
				} else {
					l.TimeoutUpdate()
					log.Println("Timeout update")
					nextDuration = duration
				}
				isWaitingForTimeout = !isWaitingForTimeout

				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(nextDuration)
			}
		}
	}()
}

func (l *Lobby) HandleJoin(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("New connection established")

	err, authMsg := l.getAuthenticationMessage(conn)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Received an octapod authentication message")

	octapod := l.identifyOctapod(authMsg.ID, authMsg.Password, conn)
	if octapod == nil {
		return
	}

	octapod.Run()
}

func (l *Lobby) getAuthenticationMessage(conn *websocket.Conn) (error, *AuthMessage) {
	messageType, content, err := conn.ReadMessage()
	if err != nil {
		return errors.New("Error reading authentication message:" + err.Error()), nil
	}

	if messageType != websocket.TextMessage {
		sendErrorAndClose(conn, "Authentication requires a text message with credentials.")
		return errors.New("received non-text message during authentication"), nil
	}

	var authMsg AuthMessage
	err = json.Unmarshal(content, &authMsg)
	if err != nil {
		sendErrorAndClose(conn, "Invalid authentication message format.")
		return errors.New("Error unmarshalling authentication message:" + err.Error()), nil
	}

	return nil, &authMsg
}

func (l *Lobby) identifyOctapod(id string, password string, conn *websocket.Conn) *Octapod {
	log.Println("Identifying octapod", id)

	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	_, isOctapodExisted := l.Octapods[id]
	var octapod *Octapod

	if isOctapodExisted {
		octapod = l.Octapods[id]
		isValidPassword := octapod.VerifyPassword(password)

		if !isValidPassword {
			log.Println("Octapod exists but password is not correct")
			sendErrorAndClose(conn, "Octapod exists but password is not correct")
			octapod.Disconnect()
			return nil
		}

		if octapod.Conn != nil {
			log.Println("Attempting to connect to already connected octapod")
			sendErrorAndClose(conn, "Octapod is already connected")
			octapod.Disconnect()
			return nil
		}

		octapod.Connect(conn)
		log.Println("Octapod [", id, "] reconnected")
	} else {
		octapod = NewOctapod(id, password, conn, l.Maze)

		l.Octapods[id] = octapod
		log.Println("New octapod [", id, "] registered")
	}

	return octapod
}

func (l *Lobby) Update() {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	for _, octapod := range l.Octapods {
		if octapod.Conn == nil {
			continue
		}
		octapod.InactiveCount++
		sensor := l.Maze.GetSensor(octapod.Position)
		octapod.Sensor <- sensor
		log.Println("Sensor data sent to octapod [", octapod.Id, "]")
	}
}

func (l *Lobby) TimeoutUpdate() {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	for _, octapod := range l.Octapods {
		if octapod.Conn == nil {
			continue
		}
		octapod.Sensor <- nil
		log.Println("Timeout signal sent to octapod", octapod.Id)
		if octapod.InactiveCount >= MaxInactive {
			octapod.Disconnect()
			log.Println("Octapod [", octapod.Id, "] disconnected due to inactivity")
		}
	}
}

func sendErrorAndClose(conn *websocket.Conn, message string) {
	errorMessage := ErrorMessage{Error: message}
	errorJson, _ := json.Marshal(errorMessage)
	conn.WriteMessage(websocket.TextMessage, errorJson)
	conn.Close()
}
