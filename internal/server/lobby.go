package server

import (
	"gbccsclub/octopod-challenge/internal/model"
	"gbccsclub/octopod-challenge/pkg"
	"github.com/google/uuid"
	"log"
	"strconv"
	"sync"
	"time"
)

type Lobby struct {
	mu         sync.Mutex
	maze       *Maze
	discordBot *DiscordBot

	ticker    *time.Ticker
	done      chan struct{}
	restart   chan struct{}
	stepCount int

	config         *Config
	AdminHandler   *AdminHandler
	OctapodHandler *OctapodHandler
	stage          model.Status
}

func NewLobby(config *Config) *Lobby {
	return &Lobby{
		stepCount:      0,
		config:         config,
		restart:        make(chan struct{}, 1),
		AdminHandler:   NewAdminHandler(config),
		OctapodHandler: NewOctapodHandler(),
		stage:          model.Exploring,
	}
}

func (l *Lobby) Start() {
	l.setupLobbyFromConfig()
	go l.Loop()
}

func (l *Lobby) setupLobbyFromConfig() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.stage = model.Exploring
	l.stepCount = 0
	l.ticker = time.NewTicker(time.Duration(l.config.TickInterval) * time.Millisecond)
	l.maze = NewMaze(l.config.MazeSize, l.config.MazeSize)
	l.maze.Generate()
	l.discordBot = NewDiscordBot(l.config.DiscordBotToken, l.config.DiscordChannelId)
	l.done = make(chan struct{})
}

func (l *Lobby) Loop() {
	for {
		select {
		case <-l.ticker.C:
			l.tick()
		case <-l.done:
			return
		case <-l.restart:
			l.handleRestart()
		}
	}
}

func (l *Lobby) handleRestart() {
	l.ticker.Stop()
	l.setupLobbyFromConfig()
}

func (l *Lobby) Stop() {
	close(l.done)
	l.ticker.Stop()
}

func (l *Lobby) Restart() {
	l.Stop()
	l.Start()
}

func (l *Lobby) RequestRestart() {
	log.Println("Restart requested")
	select {
	case l.restart <- struct{}{}:
	default:
	}
}

func (l *Lobby) renderMazeAscii() string {
	octapodPositions := l.OctapodHandler.GetOctapodPositionSet()
	view := ""
	for y := -1; y <= l.maze.Height; y++ {
		for x := -1; x <= l.maze.Width; x++ {
			pos := pkg.Vec2(x, y)

			if x == l.maze.Width-1 && y == l.maze.Height-1 {
				view += "* "
			} else if octId, ok := octapodPositions[pos]; ok {
				view += octId[0:1] + " "
			} else if l.maze.IsAvailable(pos) {
				view += "  "
			} else {
				view += "▦ "
			}
		}

		view += "\n"
	}
	return view
}

func (l *Lobby) tick() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.OctapodHandler.GetOctapodCount() == 0 {
		return
	}

	// Update octapods
	solvedOctapods := l.OctapodHandler.UpdateAll(l.maze)

	// Render maze
	view := l.renderStats(solvedOctapods)
	view += "```" + l.renderMazeAscii() + "```"
	//for _, v := range splitByNewline(view) {
	//	l.discordBot.SendMessage(v)
	//}
	l.discordBot.SendMessage(view)

	// Ping octapods
	tickId := uuid.New().String()
	l.OctapodHandler.PingAll(tickId, model.Exploring, l.maze)

	// Update step count
	l.updateStep()
}

func (l *Lobby) updateStep() {
	l.stepCount++
	if l.stage == model.Exploring && l.stepCount >= l.config.MaxExplorationSteps {
		l.stage = model.Solving
		l.stepCount = 0
	} else if l.stage == model.Solving && l.stepCount >= l.config.MaxSolvingSteps {
		l.stage = model.Ended
		l.stepCount = 0
	} else if l.stage == model.Ended {
		l.stage = model.Exploring
		l.stepCount = 0
	}
}

func (l *Lobby) renderStats(solvedOctapods []*model.Octapod) string {
	view := ""

	view += "Stage: " + l.stage.String() + "\n"

	if l.stage == model.Exploring {
		view += "Step: " + strconv.Itoa(l.stepCount) + "/" + strconv.Itoa(l.config.MaxExplorationSteps) + "\n"
	} else if l.stage == model.Solving {
		view += "Step: " + strconv.Itoa(l.stepCount) + "/" + strconv.Itoa(l.config.MaxSolvingSteps) + "\n"
	}

	view += "Octapods: " + strconv.Itoa(l.OctapodHandler.GetOctapodCount()) + "\n"

	// Display solved octapods
	if len(solvedOctapods) > 0 {
		view += "Solved: "
		for _, octapod := range solvedOctapods {
			view += octapod.GetId() + ", "
		}
		view += "\n"
	}
	return view
}

// splitByNewline splits a string to keep it under 2000 characters by newline
//func splitByNewline(s string) []string {
//	if len(s) <= 2000 {
//		return []string{s}
//	}
//
//	// Split by newline
//	lines := strings.Split(s, "\n")
//	result := []string{}
//	currentLine := ""
//
//	for _, line := range lines {
//		discordLimit := 2000
//		if len(currentLine)+len(line) > discordLimit-2 {
//			result = append(result, currentLine)
//			currentLine = ""
//		}
//		currentLine += line + "\n"
//	}
//	if currentLine != "" {
//		result = append(result, currentLine)
//	}
//	return result
//}
