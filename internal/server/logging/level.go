package logging

type Level int

const (
	LevelNoop Level = iota
	LevelDev
	LevelProd
)
