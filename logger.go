package ade_linter

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Logger struct {
	stackEntries []StackEntry
}

type StackEntry struct {
	kind string
	name string
}

func NewLogger() Logger {
	return Logger{
		stackEntries: make([]StackEntry, 0),
	}
}

func (logger Logger) AddEntry(kind string, name string) Logger {
	newStackEntries := make([]StackEntry, len(logger.stackEntries))
	copy(newStackEntries, logger.stackEntries)

	return Logger{
		stackEntries: append(newStackEntries, StackEntry{
			kind: kind,
			name: name,
		}),
	}
}

func (logger Logger) Error(err error) {
	errorMsg := fmt.Sprintf("%s %s", getErrorTag(), err.Error())
	stack := logger.getStack()

	completeMsg := append([]string{errorMsg}, stack...)

	fmt.Println(strings.Join(completeMsg, "\n"))
}

func (logger Logger) Pass(msg string) {
	passMsg := fmt.Sprintf("%s %s", getPassTag(), msg)
	stack := logger.getStack()

	completeMsg := append([]string{passMsg}, stack...)

	fmt.Println(strings.Join(completeMsg, "\n"))
}

func (logger Logger) getStack() []string {
	return Map(Reverse(logger.stackEntries), func(item StackEntry) string {
		return fmt.Sprintf("%s %s", color.HiMagentaString(fmt.Sprint("@", item.kind)), item.name)
	})
}

func getPassTag() string {
	return getLevelTag("PASS", color.New(color.FgGreen))

}

func getErrorTag() string {
	return getLevelTag("ERROR", color.New(color.FgRed))
}

func getLevelTag(level string, tone *color.Color) string {
	return fmt.Sprint(color.WhiteString("["), tone.Sprint(level), color.WhiteString("]"))
}
