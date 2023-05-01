package ade_linter

var Log = NewLogger()

func Lint(doc Document) bool {
	return areTitlesCorrect(doc) && checkGlobalInfo(doc.Sheets[GLOBAL_INFO]) && checkBoards(getBoardSheets(doc))
}

func getBoardSheets(doc Document) map[string]Sheet {
	boards := NewMap(doc.Sheets)
	delete(boards, GLOBAL_INFO)
	return boards
}
