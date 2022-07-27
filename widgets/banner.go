package widgets

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

// BannerWidget is a widget that displays the host banner
func BannerWidget(v *viper.Viper, f formatFn) (string, error) {
	v.SetDefault("banner.file", path.Join(v.GetString("config.path"), "banner.txt"))

	bannerFile := v.GetString("banner.file")
	content, err := os.ReadFile(bannerFile)

	if err != nil {
		return os.Getenv("HOSTNAME"), nil
	}

	return string(content), nil
}
