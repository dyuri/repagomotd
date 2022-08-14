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

func gradStyle(grad colorgrad.Gradient, phase float64) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(grad.At(phase).Hex()))
}

func bgGradStyle(bgrad, fgrad colorgrad.Gradient, phase float64) lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color(bgrad.At(phase).Hex())).Foreground(lipgloss.Color(fgrad.At(phase).Hex()))
}

func PBarGradient(percentage float64, width int, activeGrad, inactiveGrad colorgrad.Gradient, ltext, rtext string) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}
	pb := ""
	ch := ""
	lchars := []rune(ltext)
	rchars := []rune(rtext)
	blength := int(percentage / 100 * float64(width))
	for i := 0; i < blength; i++ {
		if i < len(lchars) {
			ch = string(lchars[i])
		} else if i > width-len(rchars) {
			ch = string(rchars[i-width+len(rchars)])
		} else {
			ch = " "
		}
		pb += bgGradStyle(activeGrad, inactiveGrad, float64(i)/float64(width)).Render(ch)
	}
	for i := blength; i < width; i++ {
		if i > width-len(rchars) {
			ch = string(rchars[i-width+len(rchars)])
		} else if i < len(lchars) {
			ch = string(lchars[i])
		} else {
			ch = " "
		}
		pb += bgGradStyle(inactiveGrad, activeGrad, float64(i)/float64(width)).Render(ch)
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
		content = lipgloss.PlaceHorizontal(lipgloss.Width(content), lipgloss.Left, content)
		content = lipgloss.PlaceHorizontal(widgetWidth-2*widgetPadding, lipgloss.Center, content)
	}

	return boxStyle.Render(content)
}

// TODO display widget name in top border
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

// TODO display widget name in top border
func BorderGradient(contents []string, grad colorgrad.Gradient) string {
	content := strings.Builder{}
	boxStyle := lipgloss.NewStyle()
	widgetWidth := GetWidgetWidth()
	boxStyle = boxStyle.MaxWidth(widgetWidth - 2)

	// calculate length of gradient
	glength := float64(widgetWidth)
	for _, c := range contents {
		glength += float64(strings.Count(c, "\n")) + 1 // +1 for last newline
	}

	ch := 0

	// first line
	content.WriteString(gradStyle(grad, 0).Render(normalBorder.TopLeft))
	for i := 0; i < widgetWidth-2; i++ {
		content.WriteString(gradStyle(grad, float64(i+1)/glength).Render(normalBorder.Horizontal))
	}
	content.WriteString(gradStyle(grad, float64(widgetWidth)/glength).Render(normalBorder.TopRight))
	content.WriteString("\n")

	for i, c := range contents {
		if i != 0 {
			// middle line
			content.WriteString(gradStyle(grad, float64(ch+1)/glength).Render(normalBorder.VerticalLeft))
			for i := 0; i < widgetWidth-2; i++ {
				content.WriteString(gradStyle(grad, float64(ch+i+2)/glength).Render(normalBorder.Horizontal))
			}
			content.WriteString(gradStyle(grad, float64(ch+widgetWidth-1)/glength).Render(normalBorder.VerticalRight))
			content.WriteString("\n")
			ch += 1
		}

		for _, line := range strings.Split(strings.TrimRight(boxStyle.Render(c), "\n"), "\n") {
			content.WriteString(gradStyle(grad, float64(ch+1)/glength).Render(normalBorder.Vertical))
			content.WriteString(line)
			content.WriteString(gradStyle(grad, float64(ch+widgetWidth-1)/glength).Render(normalBorder.Vertical))
			content.WriteString("\n")
			ch += 1
		}
	}

	// last line
	content.WriteString(gradStyle(grad, float64(ch)/glength).Render(normalBorder.BottomLeft))
	for i := 0; i < widgetWidth-2; i++ {
		content.WriteString(gradStyle(grad, float64(ch+i+1)/glength).Render(normalBorder.Horizontal))
	}
	content.WriteString(gradStyle(grad, 1).Render(normalBorder.BottomRight))
	content.WriteString("\n")

	return content.String()
}
