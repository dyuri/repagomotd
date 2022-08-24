package widgets

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/docker"
	"github.com/spf13/viper"
)

func DockerWidget(v *viper.Viper, f formatFn) (WidgetResponse, error) {
	sb := strings.Builder{}

	f1 := f("10", "", true)
	f2 := f("1", "", true)
	f3 := f("7", "", false)

	cgroups, _ := docker.GetDockerStat()

	for _, cgroup := range cgroups {
		formatter := f1
		if !cgroup.Running {
			formatter = f2
		}
		fmt.Fprintf(&sb, "%s %s\n", formatter(cgroup.Name), f3(cgroup.Status))
	}

	return WidgetResponse{
		"docker",
		sb.String(),
		"",
	}, nil
}
