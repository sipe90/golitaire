package games

type Game interface {
	Deal(number int)

	// Navigation
	Up()
	Down()
	Left()
	Right()

	Select()

	Debug()
}
