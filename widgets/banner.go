package widgets

import (
	"os"
	"path"

	"github.com/spf13/viper"
	"github.com/zcalusic/sysinfo"
)

// BannerWidget is a widget that displays the host banner
func BannerWidget(v *viper.Viper, f formatFn) (WidgetResponse, error) {
	v.SetDefault("banner.file", path.Join(v.GetString("config.path"), "banner.txt"))

	bannerFile := v.GetString("banner.file")
	content, err := os.ReadFile(bannerFile)

	if err != nil {
		var si sysinfo.SysInfo
		si.GetSysInfo()
		return WidgetResponse{
			"",
			f("1", "", true)(si.Node.Hostname),
			"center",
		}, nil
	}

	return WidgetResponse{
		"",
		string(content),
		"center",
	}, nil
}
