package handlers

type Database interface {
	Ping() error
}
