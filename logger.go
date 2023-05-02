package ade_linter

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Logger struct {
	stackEntries []StackEntry
	errorTag     string
	warnTag      string
	passTag      string
}

type StackEntry struct {
	kind string
	name string
}

func NewLogger() Logger {
	return Logger{
		stackEntries: make([]StackEntry, 0),
		errorTag:     getErrorTag(),
		warnTag:      getWarnTag(),
		passTag:      getPassTag(),
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
		errorTag: getErrorTag(),
		warnTag:  getWarnTag(),
		passTag:  getPassTag(),
	}
}

func (logger Logger) Error(err error) {
	logger.printErr(logger.errorTag, err)
}

func (logger Logger) Warn(err error) {
	logger.printErr(logger.warnTag, err)
}

func (logger Logger) Pass(msg string) {
	passMsg := fmt.Sprintf("%s %s", logger.passTag, msg)
	stack := logger.getStack()

	completeMsg := append([]string{passMsg}, stack...)

	fmt.Println(strings.Join(completeMsg, "\n"))
}

func (logger Logger) printErr(tag string, err error) {
	msg := fmt.Sprintf("%s %s", tag, err.Error())
	stack := logger.getStack()

	completeMsg := append([]string{msg}, stack...)

	fmt.Println(strings.Join(completeMsg, "\n"))
	fmt.Print("\n")
}

func (logger Logger) getStack() []string {
	return Map(Reverse(logger.stackEntries), func(item StackEntry) string {
		return fmt.Sprintf("%s %s", color.HiMagentaString(fmt.Sprint("@", item.kind)), item.name)
	})
}

func getErrorTag() string {
	return getLevelTag("ERROR", color.New(color.FgRed))
}

func getWarnTag() string {
	return getLevelTag("WARN", color.New(color.FgYellow))
}

func getPassTag() string {
	return getLevelTag("PASS", color.New(color.FgGreen))

}

func getLevelTag(level string, tone *color.Color) string {
	return fmt.Sprint(color.WhiteString("["), tone.Sprint(level), color.WhiteString("]"))
}
