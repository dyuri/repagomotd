package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/lipgloss"
	"github.com/dyuri/go-repamotd/widgets"
	"github.com/spf13/viper"
)

// WIDGETS contains the available widgets
var WIDGETS = map[string]widgets.WidgetFn{
	"naptar":  widgets.NaptarWidget,
	"banner":  widgets.BannerWidget,
	"sysinfo": widgets.SysinfoWidget,
}

// TODO parallelize
func renderWidgets(v *viper.Viper) {
	configuredWidgets := v.GetStringSlice("widgets")
	widgetContents := make([]string, 0, len(configuredWidgets))
	for _, widget := range configuredWidgets {
		if fn, ok := WIDGETS[widget]; ok {
			output, err := fn(v, widgets.Formatter)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error rendering widget: %s - %s\n", widget, err)
			} else {
				widgetContents = append(widgetContents, widgets.AlignContent(output))
			}
		}
	}
	fmt.Println(widgets.Border(widgetContents, lipgloss.NewStyle().Foreground(lipgloss.Color("4"))))
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
