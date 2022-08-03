package widgets

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mazznoer/colorgrad"
	"github.com/spf13/viper"
)

type formatFn func(string, string, bool) func(string) string
type WidgetFn func(*viper.Viper, formatFn) (WidgetResponse, error)

type WidgetResponse struct {
	Name    string
	Content string
	Place   string
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

// TODO Border, BorderColor, BorderGradient
