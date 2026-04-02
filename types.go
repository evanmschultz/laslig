package laslig

// TableWrapMode controls how long framed structured text is compacted in
// constrained widths.
type TableWrapMode string

const (
	// TableWrapAuto wraps where possible and rebalances structured content to
	// fit the available width.
	TableWrapAuto TableWrapMode = "auto"
	// TableWrapNever keeps each logical line unwrapped and truncates if needed.
	// Today this matches TableWrapTruncate intentionally; the separate name keeps
	// the API semantics explicit for callers that want to state "do not wrap".
	TableWrapNever TableWrapMode = "never"
	// TableWrapTruncate truncates long values without wrapping. Today this
	// behaves the same as TableWrapNever and differs mainly in caller intent.
	TableWrapTruncate TableWrapMode = "truncate"
)

func (mode TableWrapMode) normalized() TableWrapMode {
	switch mode {
	case TableWrapNever, TableWrapTruncate:
		return mode
	default:
		return TableWrapAuto
	}
}

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
	// MaxWidth clamps styled table width (content + frame) when rendering in
	// human mode. When omitted, the table shrinks toward content width, stays
	// within the available terminal width, and uses Läslig's readable default
	// cap.
	MaxWidth int `json:"maxWidth,omitempty"`
	// WrapMode controls how long table content is compacted in constrained
	// widths.
	WrapMode TableWrapMode `json:"wrapMode,omitempty"`
}

// Panel describes one titled boxed block.
type Panel struct {
	Title  string `json:"title,omitempty"`
	Body   string `json:"body"`
	Footer string `json:"footer,omitempty"`
	// MaxWidth caps the total panel width (content + border/padding) in human
	// mode. When omitted, the panel shrinks toward content width, stays within
	// the available terminal width, and uses Läslig's readable default cap.
	MaxWidth int `json:"maxWidth,omitempty"`
	// WrapMode controls how long panel body/footer lines are compacted.
	WrapMode TableWrapMode `json:"wrapMode,omitempty"`
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
	// MaxWidth caps the total frame width (content + frame) in styled human
	// mode. When omitted, the block shrinks toward content width, stays within
	// the available terminal width, and uses Läslig's readable default cap.
	MaxWidth int `json:"maxWidth,omitempty"`
	// WrapMode controls how long block text is compacted when constrained.
	WrapMode TableWrapMode `json:"wrapMode,omitempty"`
}

// LogBlock describes one titled boxed transcript or log excerpt.
type LogBlock struct {
	Title  string `json:"title,omitempty"`
	Body   string `json:"body"`
	Footer string `json:"footer,omitempty"`
	// MaxWidth caps the total frame width (content + frame) in styled human
	// mode. When omitted, the block shrinks toward content width, stays within
	// the available terminal width, and uses Läslig's readable default cap.
	MaxWidth int `json:"maxWidth,omitempty"`
	// WrapMode controls how long block text is compacted when constrained.
	WrapMode TableWrapMode `json:"wrapMode,omitempty"`
}
