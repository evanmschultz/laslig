package laslig

// NoticeLevel identifies one user-facing diagnostic level.
type NoticeLevel string

const (
	// NoticeInfoLevel identifies informational notices.
	NoticeInfoLevel NoticeLevel = "info"
	// NoticeSuccessLevel identifies success notices.
	NoticeSuccessLevel NoticeLevel = "success"
	// NoticeWarningLevel identifies warning notices.
	NoticeWarningLevel NoticeLevel = "warning"
	// NoticeErrorLevel identifies error notices.
	NoticeErrorLevel NoticeLevel = "error"
)

// Notice describes one user-facing diagnostic block.
type Notice struct {
	Level  NoticeLevel `json:"level"`
	Title  string      `json:"title,omitempty"`
	Body   string      `json:"body,omitempty"`
	Detail []string    `json:"detail,omitempty"`
}

// Field describes one labeled value in records and list items.
type Field struct {
	Label      string `json:"label"`
	Value      string `json:"value"`
	Identifier bool   `json:"identifier,omitempty"`
	Muted      bool   `json:"muted,omitempty"`
	Badge      bool   `json:"badge,omitempty"`
}

// Record describes one labeled data block.
type Record struct {
	Title  string  `json:"title"`
	Fields []Field `json:"fields,omitempty"`
}

// KV describes one aligned key-value block.
type KV struct {
	Title string  `json:"title,omitempty"`
	Pairs []Field `json:"pairs,omitempty"`
	Empty string  `json:"empty,omitempty"`
}

// Paragraph describes one wrapped long-form text block.
type Paragraph struct {
	Title  string `json:"title,omitempty"`
	Body   string `json:"body"`
	Footer string `json:"footer,omitempty"`
}

// ListItem describes one item in a rendered list.
type ListItem struct {
	Title  string  `json:"title"`
	Badge  string  `json:"badge,omitempty"`
	Fields []Field `json:"fields,omitempty"`
}

// List describes one titled list block.
type List struct {
	Title string     `json:"title"`
	Items []ListItem `json:"items,omitempty"`
	Empty string     `json:"empty,omitempty"`
}

// Table describes one titled table block.
type Table struct {
	Title   string     `json:"title"`
	Header  []string   `json:"header,omitempty"`
	Rows    [][]string `json:"rows,omitempty"`
	Caption string     `json:"caption,omitempty"`
	Empty   string     `json:"empty,omitempty"`
}

// Panel describes one titled boxed block.
type Panel struct {
	Title  string `json:"title,omitempty"`
	Body   string `json:"body"`
	Footer string `json:"footer,omitempty"`
}

// StatusLine describes one compact semantic status row.
type StatusLine struct {
	Level  NoticeLevel `json:"level,omitempty"`
	Label  string      `json:"label,omitempty"`
	Text   string      `json:"text"`
	Detail string      `json:"detail,omitempty"`
}

// Markdown describes one Markdown block rendered for terminal output.
type Markdown struct {
	Title  string `json:"title,omitempty"`
	Body   string `json:"body"`
	Footer string `json:"footer,omitempty"`
}

// CodeBlock describes one titled code-style block with optional language hinting.
type CodeBlock struct {
	Title    string `json:"title,omitempty"`
	Language string `json:"language,omitempty"`
	Body     string `json:"body"`
	Footer   string `json:"footer,omitempty"`
}

// LogBlock describes one titled boxed transcript or log excerpt.
type LogBlock struct {
	Title  string `json:"title,omitempty"`
	Body   string `json:"body"`
	Footer string `json:"footer,omitempty"`
}
