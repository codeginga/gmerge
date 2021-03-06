package gmerge

// GFunc represents func to run in a goroutine
type GFunc func() error

// Merger wraps metods to handle goroutine merge
type Merger interface {
	Add(string, GFunc)
	Run() map[string]error
}
