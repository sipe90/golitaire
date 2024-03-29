package golitaire

import (
	// "math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sipe90/golitaire/internal/games"
)

type model struct {
	games.Game
}

func CreateModel() model {
	return model{
		Game: games.FreeCellGame(),
	}
}

func (m model) Init() tea.Cmd {
	// number := rand.Intn(32000)
	m.Game.Deal(10221)
	// m.Game.Debug()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	game := m.Game

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "w":
			game.Up()
		case "down", "s":
			game.Down()
		case "left", "a":
			game.Left()
		case "right", "d":
			game.Right()
		case "enter", " ":
			game.Select()
		case "r":
			game.Redeal()
		case "z":
			game.Undo()
		}
	case tea.WindowSizeMsg:
		game.Resize(msg.Width, msg.Height)
	}

	return m, nil
}

func (m model) View() string {
	return m.Game.View()
}
