package games

type Game interface {
	Deal(number int)
	Redeal()

	// Navigation
	Up()
	Down()
	Left()
	Right()

	Select()
	Undo()

	View() string

	Resize(w int, h int)

	Debug()
}
