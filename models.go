package ade_linter

func NewDocument() Document {
	return Document{
		Sheets: make(map[string]Sheet),
	}
}

type Document struct {
	Sheets map[string]Sheet
}

type Sheet = [][]Cell

type Cell = string

type Tables = map[string]Table

type Table = [][]Cell

type Test interface {
	Run() bool
}

type FalseTest struct{}

func (test FalseTest) Run() bool {
	return false
}
