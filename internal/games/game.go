package games

type Game interface {
	Deal(number int)

	// Navigation
	Up()
	Down()
	Left()
	Right()

	Select()

	View() string

	Resize(w int, h int)

	Debug()
}
