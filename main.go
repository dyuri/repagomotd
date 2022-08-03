package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dyuri/go-repamotd/widgets"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// WIDGETS contains the available widgets
var WIDGETS = map[string]widgets.WidgetFn{
	"naptar":  widgets.NaptarWidget,
	"banner":  widgets.BannerWidget,
	"sysinfo": widgets.SysinfoWidget,
}

const widgetWidth = 78
const widgetPadding = 2

// TODO parallelize
func renderWidgets(v *viper.Viper) {
	var display = strings.Builder{}
	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("4")).
		Width(widgetWidth).
		PaddingLeft(widgetPadding).
		PaddingRight(widgetPadding)

	// terminal => restrict max width
	if term.IsTerminal(int(os.Stdout.Fd())) {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			boxStyle = boxStyle.MaxWidth(w)
		} else {
			boxStyle = boxStyle.MaxWidth(80)
		}
	}

	for _, widget := range v.GetStringSlice("widgets") {
		if fn, ok := WIDGETS[widget]; ok {
			output, err := fn(v, widgets.Formatter)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error rendering widget: %s - %s\n", widget, err)
			} else {
				content := output.Content
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
				display.WriteString(boxStyle.Render(content))
				display.WriteString("\n")
			}
		}
	}
	fmt.Println(display.String())
}

func main() {
	configFile := flag.String("config", "", "config file")
	createConfig := flag.Bool("create-config", false, "create config file")

	flag.Parse()

	if len(*configFile) > 0 {
		viper.SetConfigFile(*configFile)
	} else if xdg := os.Getenv("XDG_CONFIG_HOME"); len(xdg) > 0 {
		viper.AddConfigPath(path.Join(xdg, "go-repamotd"))
		viper.SetConfigName("config")
	} else if home := os.Getenv("HOME"); len(home) > 0 {
		viper.AddConfigPath(path.Join(home, ".config", "go-repamotd"))
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	// add defaults
	viper.SetDefault("widgets", []string{
		"banner",
		"sysinfo",
		"naptar",
	})

	// root config dir
	configPath := ""
	if xdg := os.Getenv("XDG_CONFIG_HOME"); len(xdg) > 0 {
		configPath = path.Join(xdg, "go-repamotd")
	} else if home := os.Getenv("HOME"); len(home) > 0 {
		configPath = path.Join(home, ".config", "go-repamotd")
	}
	viper.SetDefault("config.path", configPath)

	// read config
	if err := viper.ReadInConfig(); err != nil && *createConfig {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("no config file found")
			if len(configPath) > 0 {
				fmt.Printf("creating default config in: %s\n", configPath)
				os.MkdirAll(configPath, 0755)
				if err := viper.WriteConfigAs(path.Join(configPath, "config.yaml")); err != nil {
					fmt.Printf("error creating default config: %s\n", err)
				}
			}
		} else {
			fmt.Println("error reading config file:", err)
		}
	}

	renderWidgets(viper.GetViper())
}
