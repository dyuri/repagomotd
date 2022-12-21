package widgets

import (
	"os"
	"path"

	"github.com/dyuri/go-figure"
	"github.com/spf13/viper"
	"github.com/zcalusic/sysinfo"
)

// TODO sysinfo => gopsutil

// BannerWidget is a widget that displays the host banner
func BannerWidget(v *viper.Viper, f formatFn) (WidgetResponse, error) {
	v.SetDefault("banner.file", path.Join(v.GetString("config.path"), "banner.txt"))
	v.SetDefault("banner.font", "3d")

	bannerFile := v.GetString("banner.file")
	content, err := os.ReadFile(bannerFile)

	if err != nil {
		var si sysinfo.SysInfo
		si.GetSysInfo()
		fig := figure.NewFigure(si.Node.Hostname, v.GetString("banner.font"), true)
		content := fig.ColorString(
			figure.GradientRGBColorizer(
				184, 187, 38,
				215, 153, 33,
			),
		)
		return WidgetResponse{
			"",
			f("", "", true)(content),
			"center",
		}, nil
	}

	return WidgetResponse{
		"",
		string(content),
		"center",
	}, nil
}
