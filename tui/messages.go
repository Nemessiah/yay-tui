package tui

// yaySearchResultMsg and yaySearchErrorMsg are async messages sent to Update().

type yaySearchResultMsg struct {
	results []string
}

type yaySearchErrorMsg struct {
	errorText string
}
