package chat

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Message struct {
	Text *string `json:"text,omitempty"`

	Color string `json:"color,omitempty"`

	Bold          bool `json:"bold,omitempty,string"`
	Italic        bool `json:"italic,omitempty,string"`
	Underlined    bool `json:"underlined,omitempty,string"`
	Strikethrough bool `json:"strikethrough,omitempty,string"`
	Obfuscated    bool `json:"obfuscated,omitempty,string"`

	Extra []Message `json:"extra,omitempty"`

	Translate string    `json:"translate,omitempty"`
	With      []Message `json:"with,omitempty"`

	ClickEvent *ClickEvent `json:"clickEvent,omitempty"`
	HoverEvent *HoverEvent `json:"hoverEvent,omitempty"`
}

type ClickEvent struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}

type HoverEvent struct {
	Action   string      `json:"action"`
	Contents interface{} `json:"value"`
}

func NewMessage(s string) (m Message) {
	var component Message
	s = strings.ReplaceAll(s, "§", "&")
	for i := 0; i < len(s); i++ {
		if s[i] == '&' {
			if i+1 == len(s) {
				break
			}
			i++

			t := s[i]
			switch {
			case t == 'r':
				component.Color = "reset"
				component.Bold = false
				component.Italic = false
				component.Underlined = false
				component.Strikethrough = false
				component.Obfuscated = false
				
			case t >= '0' && t <= '9' || t >= 'a' && t <= 'f':
				component.Color = colors(t)

			case t >= 'k' && t <= 'o':
				styles(t, &component)
			}

			continue
		}

		var n int
		for n = i; n < len(s); n++ {
			if s[n] == '&' {
				break
			}
		}

		x := s[i:n]
		component.Text = &x
		if m.Text == nil {
			m = component
		} else {
			m.Extra = append(m.Extra, component)
			component = Message{}
		}
		i = n - 1
	}

	return
}

func (m *Message) String() string {
	v, _ := json.Marshal(m)
	return string(v)
}

// Opens the url for the player
func (m Message) WithOpenURLClickEvent(url string) Message {
	m.ClickEvent = &ClickEvent{
		Action: "open_url",
		Value:  url,
	}
	return m
}

// Causes the player to send the command
// If not prefixed with '/', it will send the messsage
func (m Message) WithRunCommandClickEvent(cmd string) Message {
	m.ClickEvent = &ClickEvent{
		Action: "run_command",
		Value:  cmd,
	}
	return m
}

// Fills the player's chat with the command
// Also works with chat messages (no / prefix)
func (m Message) WithSuggestCommandClickEvent(cmd string) Message {
	m.ClickEvent = &ClickEvent{
		Action: "suggest_command",
		Value:  cmd,
	}
	return m
}

// Copies the text to the player's clipboard
func (m Message) WithCopyToClipboardClickEvent(text string) Message {
	m.ClickEvent = &ClickEvent{
		Action: "copy_to_clipboard",
		Value: text,
	}
	return m
}


// Shows the text
func (m Message) WithShowTextHoverEvent(msg Message) Message {
	m.HoverEvent = &HoverEvent{
		Action:   "show_text",
		Contents: msg,
	}
	return m
}

func (m Message) WithShowEntityHoverEvent(id string, name string, typ *string) Message {
	text := fmt.Sprintf(`{id:%s,name:%s}`, id, name)
	if typ != nil {
		text = strings.TrimPrefix(text, "}") + fmt.Sprintf(`, type:%s}`, *typ)
	}
	m.HoverEvent = &HoverEvent{
		Action:   "show_entity",
		Contents: text,
	}
	return m
}

func styles(b byte, msg *Message) bool {
	switch b {

	default:
		return false

	case 'k':
		msg.Obfuscated = true

	case 'l':
		msg.Bold = true

	case 'm':
		msg.Strikethrough = true

	case 'n':
		msg.Underlined = true

	case 'o':
		msg.Italic = true
	}

	return true
}

func colors(b byte) string {
	switch b {

	default:
		return ""

	case '0':
		return "black"

	case '1':
		return "dark_blue"

	case '2':
		return "dark_green"

	case '3':
		return "dark_aqua"

	case '4':
		return "dark_red"

	case '5':
		return "dark_purple"

	case '6':
		return "gold"

	case '7':
		return "gray"

	case '8':
		return "dark_gray"

	case '9':
		return "blue"

	case 'a':
		return "green"

	case 'b':
		return "aqua"

	case 'c':
		return "red"

	case 'd':
		return "light_purple"

	case 'e':
		return "yellow"

	case 'f':
		return "white"
	}
}

func Translate(msg string, with ...Message) Message {
	return Message{Translate: msg, With: with}
}
