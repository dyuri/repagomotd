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

func pbar(percentage float64, width int, bar, space string) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}
	blength := int(percentage / 100 * float64(width))
	return strings.Repeat(bar, blength) + strings.Repeat(space, width-blength)
}

func pbarColor(percentage float64, width int, bar, space string, activeColor, inactiveColor string) string {
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

func pbarGradient(percentage float64, width int, bar, space string, activeGrad, inactiveGrad colorgrad.Gradient) string {
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
