package codearea

// State keeps the state of the widget. Its access must be synchronized through
// the mutex.
type State struct {
	CodeBuffer  CodeBuffer
	PendingCode PendingCode
}

// CodeBuffer represents the state of the buffer.
type CodeBuffer struct {
	// Content of the buffer.
	Content string
	// Position of the dot (more commonly known as the cursor), as a byte index
	// into Content.
	Dot int
}

// PendingCode represents pending code, such as during completion.
type PendingCode struct {
	// Beginning index of the text area that the pending code replaces, as a
	// byte index into RawState.Code.
	From int
	// End index of the text area that the pending code replaces, as a byte
	// index into RawState.Code.
	To int
	// The content of the pending code.
	Content string
}

func (c *CodeBuffer) InsertAtDot(text string) {
	*c = CodeBuffer{
		Content: c.Content[:c.Dot] + text + c.Content[c.Dot:],
		Dot:     c.Dot + len(text),
	}
}