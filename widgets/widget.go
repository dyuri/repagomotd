package widgets

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mazznoer/colorgrad"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

const widgetPadding = 2

type formatFn func(string, string, bool) func(string) string
type WidgetFn func(*viper.Viper, formatFn) (WidgetResponse, error)

type WidgetResponse struct {
	Name    string
	Content string
	Place   string
}

type BorderChars struct {
	Vertical      string
	Horizontal    string
	TopLeft       string
	TopRight      string
	BottomLeft    string
	BottomRight   string
	VerticalLeft  string
	VerticalRight string
}

var normalBorder = BorderChars{
	Vertical:      "│",
	Horizontal:    "─",
	TopLeft:       "╭",
	TopRight:      "╮",
	BottomLeft:    "╰",
	BottomRight:   "╯",
	VerticalLeft:  "├",
	VerticalRight: "┤",
}

func Formatter(fg, bg string, bold bool) func(string) string {
	var style lipgloss.Style
	if fg == "" {
		fg = "7"
	}
	style = lipgloss.NewStyle().Foreground(lipgloss.Color(fg))
	if bg != "" {
		style = style.Background(lipgloss.Color(bg))
	}
	if bold {
		style = style.Bold(true)
	} else {
		style = style.Bold(false)
	}

	return style.Render
}

func PBar(percentage float64, width int, bar, space string) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}
	blength := int(percentage / 100 * float64(width))
	return strings.Repeat(bar, blength) + strings.Repeat(space, width-blength)
}

func PBarColor(percentage float64, width int, bar, space string, activeColor, inactiveColor string) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}
	blength := int(percentage / 100 * float64(width))
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(activeColor))
	inactiveStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(inactiveColor))
	return activeStyle.Render(strings.Repeat(bar, blength)) + inactiveStyle.Render(strings.Repeat(space, width-blength))
}

func PBarGradient(percentage float64, width int, bar, space string, activeGrad, inactiveGrad colorgrad.Gradient) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}
	pb := ""
	blength := int(percentage / 100 * float64(width))
	for i := 0; i < blength; i++ {
		pb += lipgloss.NewStyle().Foreground(lipgloss.Color(activeGrad.At(float64(i) / float64(width)).Hex())).Render(bar)
	}
	for i := blength; i < width; i++ {
		pb += lipgloss.NewStyle().Foreground(lipgloss.Color(inactiveGrad.At(float64(i) / float64(width)).Hex())).Render(space)
	}

	return pb
}

func GetWidgetWidth() int {
	widgetWidth := 80

	// terminal => restrict max width
	if term.IsTerminal(int(os.Stdout.Fd())) {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil && w < 80 {
			widgetWidth = w
		}
	}

	return widgetWidth
}

func AlignContent(output WidgetResponse) string {
	content := output.Content
	widgetWidth := GetWidgetWidth()
	boxStyle := lipgloss.NewStyle().
		Width(widgetWidth).
		PaddingLeft(widgetPadding).
		PaddingRight(widgetPadding)

	rcontent := []rune(content)
	length := len(rcontent)
	if length > 0 && rcontent[length-1] == '\n' {
		rcontent = rcontent[:length-1]
		content = string(rcontent)
	}

	if output.Place == "center" {
		content = lipgloss.PlaceHorizontal(lipgloss.Width(content), lipgloss.Left, lipgloss.NewStyle().Background(lipgloss.Color("0")).Render(content))
		content = lipgloss.PlaceHorizontal(widgetWidth-2*widgetPadding, lipgloss.Center, content)
	}

	return boxStyle.Render(content)
}

func Border(contents []string, borderStyle lipgloss.Style) string {
	content := strings.Builder{}
	boxStyle := lipgloss.NewStyle()
	widgetWidth := GetWidgetWidth()
	boxStyle = boxStyle.MaxWidth(widgetWidth - 2)

	// first line
	content.WriteString(borderStyle.Render(normalBorder.TopLeft + strings.Repeat(normalBorder.Horizontal, widgetWidth-2) + normalBorder.TopRight))
	content.WriteString("\n")

	for i, c := range contents {
		if i != 0 {
			// middle line
			content.WriteString(borderStyle.Render(normalBorder.VerticalLeft + strings.Repeat(normalBorder.Horizontal, widgetWidth-2) + normalBorder.VerticalRight))
			content.WriteString("\n")
		}
		for _, line := range strings.Split(strings.TrimRight(boxStyle.Render(c), "\n"), "\n") {
			content.WriteString(borderStyle.Render(normalBorder.Vertical) + line + borderStyle.Render(normalBorder.Vertical))
			content.WriteString("\n")
		}
	}

	// last line
	content.WriteString(borderStyle.Render(normalBorder.BottomLeft + strings.Repeat(normalBorder.Horizontal, widgetWidth-2) + normalBorder.BottomRight))
	content.WriteString("\n")

	return content.String()
}

// TODO BorderGradient
