package widgets

import (
	"github.com/spf13/viper"
)

type formatFn func(string, string, bool) func(string) string
type WidgetFn func(*viper.Viper, formatFn) (WidgetResponse, error)

type WidgetResponse struct {
	Name    string
	Content string
	Place   string
}
