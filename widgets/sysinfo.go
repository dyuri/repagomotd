package widgets

import (
	"fmt"
	"strings"
	"time"

	"github.com/mazznoer/colorgrad"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/viper"
)

// SysinfoWidget is a widget that displays the host banner
func SysinfoWidget(v *viper.Viper, f formatFn) (WidgetResponse, error) {
	// formatters
	f1 := f("7", "0", false)
	f2 := f("11", "0", false)
	f3 := f("2", "0", false)
	f4 := f("10", "0", true)
	sb := strings.Builder{}

	addLine := func(title, value string) {
		fmt.Fprintf(&sb, "%s %s\n", f1(fmt.Sprintf("%8s", title)), value)
	}

	hostinfo, _ := host.Info()

	// node
	addLine("host:", f2(hostinfo.Hostname))

	// os
	addLine("os:", f2(hostinfo.Platform)+f3(" ["+hostinfo.PlatformVersion+"]"))

	// kernel
	addLine("kernel:", f2(hostinfo.KernelVersion)+f3(" ["+hostinfo.KernelArch+"]"))

	// CPU
	cpuInfo, _ := cpu.Info()

	cpus := map[string]cpu.InfoStat{}
	cores := map[string]map[string]struct{}{}
	threads := map[string]int{}

	// count cores and threads
	for _, c := range cpuInfo {
		phid := c.PhysicalID
		cid := c.CoreID

		// add to physical cpus
		if _, ok := cpus[phid]; !ok {
			cpus[phid] = c
		}

		// add to cores
		if _, ok := cores[phid]; !ok {
			cores[phid] = make(map[string]struct{})
		}
		cores[phid][cid] = struct{}{}

		if _, ok := threads[phid]; !ok {
			threads[phid] = 0
		}
		threads[phid]++
	}

	for i, c := range cpus {
		title := ""
		if i == "0" {
			title = "cpu:"
		}
		addLine(title, f2(c.ModelName)+f3(fmt.Sprintf(" [%d/%d]", len(cores[i]), threads[i])))
	}

	// uptime
	addLine("uptime:", f2((time.Duration(hostinfo.Uptime)*time.Second).String())+f3(fmt.Sprintf(" [%d processes]", hostinfo.Procs)))

	// load
	loadInfo, _ := load.Avg()
	addLine("load:", f4(fmt.Sprintf("%.2f", loadInfo.Load1))+f3(" (1m) ")+f4(fmt.Sprintf("%.2f", loadInfo.Load5))+f3(" (5m) ")+f4(fmt.Sprintf("%.2f", loadInfo.Load15))+f3(" (15m)"))

	// Memory
	memory, err := mem.VirtualMemory()
	memValue := f1("???")

	if err == nil {
		memValue = f4(fmt.Sprintf("%.2f", float64(memory.Used)/(1024*1024*1024))) + f2(fmt.Sprintf("/%.2fGi", float64(memory.Total)/(1024*1024*1024)))
	}
	addLine("memory:", memValue)

	grad, _ := colorgrad.NewGradient().HtmlColors("#b8bb26", "#fabd2f", "#fb4934").Build()
	grad2, _ := colorgrad.NewGradient().HtmlColors("#98971a", "#d79921", "#cc241d").Build()
	addLine("", PBarGradient(memory.UsedPercent, 64, "█", "░", grad, grad2))

	return WidgetResponse{
		"sysinfo",
		sb.String(),
		"",
	}, nil
}
