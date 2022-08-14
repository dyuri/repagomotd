package widgets

import (
	"fmt"
	"strings"

	"github.com/mazznoer/colorgrad"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/spf13/viper"
)

/****
 * DiskWidget is a widget that displays partition usage statistics
 * Configuration options:
 *   - disk.partitions: list of disks to display [optional]
 */
func DiskWidget(v *viper.Viper, f formatFn) (WidgetResponse, error) {
	sb := strings.Builder{}

	grad, _ := colorgrad.NewGradient().HtmlColors("#b8bb26", "#fabd2f", "#fb4934").Build()
	grad2, _ := colorgrad.NewGradient().HtmlColors("#484d00", "#5e4e00", "#500000").Build()

	configuredPartitions := v.GetStringSlice("disk.partitions")
	fmt.Println(configuredPartitions)

	var mps []string
	if len(configuredPartitions) == 0 {
		partitions, _ := disk.Partitions(false)
		for _, p := range partitions {
			mps = append(mps, p.Mountpoint)
		}
	} else {
		mps = configuredPartitions
	}
	for _, partition := range mps {
		usage, _ := disk.Usage(partition)
		if usage != nil {
			fmt.Fprintf(&sb, "%s\n", PBarGradient(usage.UsedPercent, GetWidgetWidth()-2*widgetPadding-2, grad, grad2, fmt.Sprintf(" %s", partition), fmt.Sprintf(" %.0f%% ", usage.UsedPercent)))
		} else {
			fmt.Fprintf(&sb, "[!] error checking partition '%s'\n", partition)
		}
	}
	return WidgetResponse{
		"disk",
		sb.String(),
		"",
	}, nil
}
