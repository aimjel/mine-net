package chat

type Message struct {
	Text string `json:"text"`

	Color string `json:"color,omitempty"`

	Bold          bool `json:"bold,omitempty,string"`
	Italic        bool `json:"italic,omitempty,string"`
	Underlined    bool `json:"underlined,omitempty,string"`
	Strikethrough bool `json:"strikethrough,omitempty,string"`
	Obfuscated    bool `json:"obfuscated,omitempty,string"`

	Extra []Message `json:"extra,omitempty"`
}

func NewMessage(s string) (m Message) {

	var component Message
	for i := 0; i < len(s); i++ {
		if s[i] == '&' {
			if i+1 == len(s) {
				break
			}
			i++

			t := s[i]
			switch {

			case t >= '0' && t <= '9' || t >= 'a' && t <= 'f' || t == 'r':
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

		component.Text = s[i:n]
		if m.Text == "" {
			m = component
		} else {
			m.Extra = append(m.Extra, component)
			component = Message{}
		}
		i = n - 1
	}

	return
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

	case 'r':
		return "reset"
	}
}
