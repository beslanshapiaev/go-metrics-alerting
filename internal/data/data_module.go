package data

type DataModule interface {
	Ping() bool
	Close()
}
