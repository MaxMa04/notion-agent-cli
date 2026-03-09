package cmd

import (
	"testing"
)

func TestParseMarkdownToBlocks(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
		checkFirst func(t *testing.T, block map[string]interface{})
	}{
		{
			name:      "heading 1",
			input:     "# Hello",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "heading_1" {
					t.Errorf("type = %v, want heading_1", b["type"])
				}
			},
		},
		{
			name:      "heading 2",
			input:     "## Sub heading",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "heading_2" {
					t.Errorf("type = %v, want heading_2", b["type"])
				}
			},
		},
		{
			name:      "heading 3",
			input:     "### Sub sub heading",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "heading_3" {
					t.Errorf("type = %v, want heading_3", b["type"])
				}
			},
		},
		{
			name:      "bullet list",
			input:     "- item one\n- item two\n- item three",
			wantCount: 3,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "bulleted_list_item" {
					t.Errorf("type = %v, want bulleted_list_item", b["type"])
				}
			},
		},
		{
			name:      "bullet with asterisk",
			input:     "* item",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "bulleted_list_item" {
					t.Errorf("type = %v, want bulleted_list_item", b["type"])
				}
			},
		},
		{
			name:      "numbered list",
			input:     "1. first\n2. second",
			wantCount: 2,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "numbered_list_item" {
					t.Errorf("type = %v, want numbered_list_item", b["type"])
				}
			},
		},
		{
			name:      "quote",
			input:     "> This is a quote",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "quote" {
					t.Errorf("type = %v, want quote", b["type"])
				}
			},
		},
		{
			name:      "divider",
			input:     "---",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "divider" {
					t.Errorf("type = %v, want divider", b["type"])
				}
			},
		},
		{
			name:      "code block",
			input:     "```go\nfmt.Println(\"hello\")\n```",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "code" {
					t.Errorf("type = %v, want code", b["type"])
				}
				code, ok := b["code"].(map[string]interface{})
				if !ok {
					t.Fatal("missing code block data")
				}
				if code["language"] != "go" {
					t.Errorf("language = %v, want go", code["language"])
				}
			},
		},
		{
			name:      "code block no language",
			input:     "```\nsome code\n```",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				code := b["code"].(map[string]interface{})
				if code["language"] != "plain text" {
					t.Errorf("language = %v, want 'plain text'", code["language"])
				}
			},
		},
		{
			name:      "todo unchecked",
			input:     "- [ ] do this",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "to_do" {
					t.Errorf("type = %v, want to_do", b["type"])
				}
				td := b["to_do"].(map[string]interface{})
				if td["checked"] != false {
					t.Error("checked should be false")
				}
			},
		},
		{
			name:      "todo checked",
			input:     "- [x] done",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				td := b["to_do"].(map[string]interface{})
				if td["checked"] != true {
					t.Error("checked should be true")
				}
			},
		},
		{
			name:      "paragraph fallback",
			input:     "Just a regular paragraph",
			wantCount: 1,
			checkFirst: func(t *testing.T, b map[string]interface{}) {
				if b["type"] != "paragraph" {
					t.Errorf("type = %v, want paragraph", b["type"])
				}
			},
		},
		{
			name:      "empty lines skipped",
			input:     "\n\n\nHello\n\n\n",
			wantCount: 1,
		},
		{
			name:      "mixed content",
			input:     "# Title\n\nA paragraph.\n\n- bullet one\n- bullet two\n\n> a quote\n\n---",
			wantCount: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := parseMarkdownToBlocks(tt.input)
			if len(blocks) != tt.wantCount {
				t.Errorf("got %d blocks, want %d", len(blocks), tt.wantCount)
				for i, b := range blocks {
					t.Logf("  block[%d]: type=%v", i, b["type"])
				}
				return
			}
			if tt.checkFirst != nil && len(blocks) > 0 {
				tt.checkFirst(t, blocks[0])
			}
		})
	}
}

func TestMakeTextBlock(t *testing.T) {
	block := makeTextBlock("paragraph", "Hello World")
	if block["type"] != "paragraph" {
		t.Errorf("type = %v, want paragraph", block["type"])
	}
	if block["object"] != "block" {
		t.Errorf("object = %v, want block", block["object"])
	}
	p, ok := block["paragraph"].(map[string]interface{})
	if !ok {
		t.Fatal("missing paragraph data")
	}
	rt, ok := p["rich_text"].([]map[string]interface{})
	if !ok || len(rt) != 1 {
		t.Fatal("expected 1 rich_text element")
	}
	text := rt[0]["text"].(map[string]interface{})
	if text["content"] != "Hello World" {
		t.Errorf("content = %v, want 'Hello World'", text["content"])
	}
}

func TestParseCSVRow(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{`"hello, world",b,c`, []string{"hello, world", "b", "c"}},
		{"single", []string{"single"}},
		{"a,,c", []string{"a", "", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseCSVRow(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d fields, want %d: %v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("field[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuildTableRow(t *testing.T) {
	row := buildTableRow([]string{"A", "B", "C"})

	if row["type"] != "table_row" {
		t.Fatalf("type = %v, want table_row", row["type"])
	}

	rowData, ok := row["table_row"].(map[string]interface{})
	if !ok {
		t.Fatal("missing table_row data")
	}

	cells, ok := rowData["cells"].([]interface{})
	if !ok {
		t.Fatal("missing cells")
	}
	if len(cells) != 3 {
		t.Fatalf("got %d cells, want 3", len(cells))
	}

	// Check first cell
	cell0 := cells[0].([]interface{})
	rt := cell0[0].(map[string]interface{})
	text := rt["text"].(map[string]interface{})
	if text["content"] != "A" {
		t.Errorf("cell[0] content = %v, want A", text["content"])
	}
}

func TestBuildTableBlock(t *testing.T) {
	rows := [][]string{
		{"Name", "Role", "Status"},
		{"Alice", "Dev", "Active"},
		{"Bob", "PM", "Done"},
	}

	block := buildTableBlock(rows, true)

	if block["type"] != "table" {
		t.Fatalf("type = %v, want table", block["type"])
	}

	tableData, ok := block["table"].(map[string]interface{})
	if !ok {
		t.Fatal("missing table data")
	}

	if tableData["table_width"] != 3 {
		t.Errorf("table_width = %v, want 3", tableData["table_width"])
	}
	if tableData["has_column_header"] != true {
		t.Error("has_column_header should be true")
	}
	if tableData["has_row_header"] != false {
		t.Error("has_row_header should be false")
	}

	children, ok := tableData["children"].([]interface{})
	if !ok {
		t.Fatal("missing children")
	}
	if len(children) != 3 {
		t.Errorf("got %d children, want 3", len(children))
	}
}

func TestBuildTableBlock_NoHeader(t *testing.T) {
	rows := [][]string{{"a", "b"}, {"c", "d"}}
	block := buildTableBlock(rows, false)
	tableData := block["table"].(map[string]interface{})
	if tableData["has_column_header"] != false {
		t.Error("has_column_header should be false")
	}
}

func TestBuildTableBlock_UnevenRows(t *testing.T) {
	rows := [][]string{
		{"a", "b", "c"},
		{"x"},
	}
	block := buildTableBlock(rows, true)
	tableData := block["table"].(map[string]interface{})
	if tableData["table_width"] != 3 {
		t.Errorf("table_width = %v, want 3 (widest row)", tableData["table_width"])
	}
	// Second row should be padded to 3 cells
	children := tableData["children"].([]interface{})
	row2 := children[1].(map[string]interface{})
	cells := row2["table_row"].(map[string]interface{})["cells"].([]interface{})
	if len(cells) != 3 {
		t.Errorf("row2 has %d cells, want 3 (padded)", len(cells))
	}
}

func TestIsSeparatorRow(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"---|---", true},
		{" --- | --- | --- ", true},
		{":---:|:---:", true},
		{"hello | world", false},
		{"---", true},
		{"", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isSeparatorRow(tt.input); got != tt.want {
				t.Errorf("isSeparatorRow(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseTableRow(t *testing.T) {
	got := parseTableRow("| Name | Role | Status |")
	want := []string{"Name", "Role", "Status"}
	if len(got) != len(want) {
		t.Fatalf("got %d cells, want %d: %v", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("cell[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseMarkdownTable(t *testing.T) {
	input := "| Name | Role |\n| --- | --- |\n| Alice | Dev |\n| Bob | PM |"
	blocks := parseMarkdownToBlocks(input)
	if len(blocks) != 1 {
		t.Fatalf("got %d blocks, want 1", len(blocks))
	}
	if blocks[0]["type"] != "table" {
		t.Errorf("type = %v, want table", blocks[0]["type"])
	}
	tableData := blocks[0]["table"].(map[string]interface{})
	children := tableData["children"].([]interface{})
	// 3 rows: header + 2 data (separator is skipped)
	if len(children) != 3 {
		t.Errorf("got %d rows, want 3", len(children))
	}
}

func TestParseMarkdownTable_Mixed(t *testing.T) {
	input := "# Title\n\n| A | B |\n| --- | --- |\n| 1 | 2 |\n\nSome text"
	blocks := parseMarkdownToBlocks(input)
	if len(blocks) != 3 {
		t.Fatalf("got %d blocks, want 3 (heading + table + paragraph)", len(blocks))
	}
	if blocks[0]["type"] != "heading_1" {
		t.Errorf("block[0] type = %v, want heading_1", blocks[0]["type"])
	}
	if blocks[1]["type"] != "table" {
		t.Errorf("block[1] type = %v, want table", blocks[1]["type"])
	}
	if blocks[2]["type"] != "paragraph" {
		t.Errorf("block[2] type = %v, want paragraph", blocks[2]["type"])
	}
}

func TestMapBlockTypeAliases(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"h1", "heading_1"},
		{"h2", "heading_2"},
		{"h3", "heading_3"},
		{"heading1", "heading_1"},
		{"heading2", "heading_2"},
		{"heading3", "heading_3"},
		{"bullet", "bulleted_list_item"},
		{"numbered", "numbered_list_item"},
		{"todo", "to_do"},
		{"p", "paragraph"},
		{"paragraph", "paragraph"},
		{"quote", "quote"},
		{"code", "code"},
		{"callout", "callout"},
		{"divider", "divider"},
		{"table", "table"},
		{"tbl", "table"},
		{"table_row", "table_row"},
		{"trow", "table_row"},
		// passthrough for native Notion types
		{"heading_1", "heading_1"},
		{"bulleted_list_item", "bulleted_list_item"},
		{"unknown_type", "unknown_type"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := mapBlockType(tt.input)
			if got != tt.want {
				t.Errorf("mapBlockType(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
