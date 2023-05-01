package ade_linter

import (
	"fmt"
	"regexp"
)

const (
	GLOBAL_INFO = "GLOBAL INFO"
	BOARD       = "BOARD "
)

func areTitlesCorrect(doc Document) bool {
	titleLogger := Log.AddEntry("Titles", "")
	return hasGlobalInfo(doc.Sheets, titleLogger) && haveBoardPrefix(getBoardSheets(doc), titleLogger)
}

func hasGlobalInfo(sheets map[string]Sheet, logger Logger) bool {
	_, ok := sheets[GLOBAL_INFO]

	if !ok {
		logger.Error(fmt.Errorf("%s sheet not found", GLOBAL_INFO))
	}

	return ok
}

func haveBoardPrefix(sheets map[string]Sheet, logger Logger) bool {
	prefixExp := regexp.MustCompile(fmt.Sprintf(`^%s\w+$`, BOARD))

	return EveryMap(sheets, func(name string, sheets Sheet) bool {
		if !prefixExp.MatchString(name) {
			logger.Error(fmt.Errorf("sheet %s doesn't have %s prefix", name, BOARD))
			return false
		}

		return true
	})
}
