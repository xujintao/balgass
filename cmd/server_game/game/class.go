package game

type ClassChar int

const (
	Wizard ClassChar = iota
	Knight
	Elf
	Magumsa
	DarkLord
	Summoner
	RageFighter
	GrowLancer
)

const MaxClassChar int = 8
