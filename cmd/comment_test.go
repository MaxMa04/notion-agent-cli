package cmd

import (
	"testing"
)

func TestParseRichText_PlainText(t *testing.T) {
	parts := parseRichText("Hello world")
	if len(parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(parts))
	}
	p := parts[0].(map[string]interface{})
	if p["type"] != "text" {
		t.Errorf("expected type text, got %v", p["type"])
	}
	txt := p["text"].(map[string]interface{})
	if txt["content"] != "Hello world" {
		t.Errorf("expected 'Hello world', got %v", txt["content"])
	}
}

func TestParseRichText_SingleMention(t *testing.T) {
	parts := parseRichText("Hey @266cdaa8-27bf-41a2-8058-90c81a416abb bitte reviewen")
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}

	// Part 0: text "Hey "
	p0 := parts[0].(map[string]interface{})
	if p0["type"] != "text" {
		t.Errorf("part 0: expected text, got %v", p0["type"])
	}
	if p0["text"].(map[string]interface{})["content"] != "Hey " {
		t.Errorf("part 0: wrong content: %v", p0["text"])
	}

	// Part 1: mention
	p1 := parts[1].(map[string]interface{})
	if p1["type"] != "mention" {
		t.Errorf("part 1: expected mention, got %v", p1["type"])
	}
	mention := p1["mention"].(map[string]interface{})
	if mention["type"] != "user" {
		t.Errorf("part 1: expected user mention, got %v", mention["type"])
	}
	user := mention["user"].(map[string]interface{})
	if user["id"] != "266cdaa8-27bf-41a2-8058-90c81a416abb" {
		t.Errorf("part 1: wrong user ID: %v", user["id"])
	}

	// Part 2: text " bitte reviewen"
	p2 := parts[2].(map[string]interface{})
	if p2["type"] != "text" {
		t.Errorf("part 2: expected text, got %v", p2["type"])
	}
	if p2["text"].(map[string]interface{})["content"] != " bitte reviewen" {
		t.Errorf("part 2: wrong content: %v", p2["text"])
	}
}

func TestParseRichText_MentionAtStart(t *testing.T) {
	parts := parseRichText("@266cdaa8-27bf-41a2-8058-90c81a416abb done")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	if parts[0].(map[string]interface{})["type"] != "mention" {
		t.Errorf("part 0: expected mention")
	}
	if parts[1].(map[string]interface{})["type"] != "text" {
		t.Errorf("part 1: expected text")
	}
}

func TestParseRichText_MentionAtEnd(t *testing.T) {
	parts := parseRichText("Check this @266cdaa8-27bf-41a2-8058-90c81a416abb")
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	if parts[0].(map[string]interface{})["type"] != "text" {
		t.Errorf("part 0: expected text")
	}
	if parts[1].(map[string]interface{})["type"] != "mention" {
		t.Errorf("part 1: expected mention")
	}
}

func TestParseRichText_MultipleMentions(t *testing.T) {
	parts := parseRichText("@aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee und @11111111-2222-3333-4444-555555555555 check")
	if len(parts) != 4 {
		t.Fatalf("expected 4 parts, got %d", len(parts))
	}
	// mention, text " und ", mention, text " check"
	if parts[0].(map[string]interface{})["type"] != "mention" {
		t.Errorf("part 0: expected mention")
	}
	if parts[1].(map[string]interface{})["type"] != "text" {
		t.Errorf("part 1: expected text")
	}
	if parts[2].(map[string]interface{})["type"] != "mention" {
		t.Errorf("part 2: expected mention")
	}
	if parts[3].(map[string]interface{})["type"] != "text" {
		t.Errorf("part 3: expected text")
	}
}

func TestParseRichText_NoFalsePositive(t *testing.T) {
	// @Max should NOT be treated as a mention (not a UUID)
	parts := parseRichText("Hey @Max check this")
	if len(parts) != 1 {
		t.Fatalf("expected 1 part (no mention), got %d", len(parts))
	}
	if parts[0].(map[string]interface{})["type"] != "text" {
		t.Errorf("expected text, got %v", parts[0].(map[string]interface{})["type"])
	}
}
