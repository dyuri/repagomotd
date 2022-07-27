package widgets

import (
	"fmt"
	"strings"

	"github.com/mackerelio/go-osstat/memory"
	"github.com/spf13/viper"
	"github.com/zcalusic/sysinfo"
)

// SysinfoWidget is a widget that displays the host banner
func SysinfoWidget(v *viper.Viper, f formatFn) (WidgetResponse, error) {
	var si sysinfo.SysInfo
	si.GetSysInfo()

	sb := strings.Builder{}

	// node
	sb.WriteString("host: ")
	sb.WriteString(si.Node.Hostname)
	sb.WriteString("\n")

	// os
	sb.WriteString("os: ")
	sb.WriteString(si.OS.Name)
	sb.WriteString(" [")
	sb.WriteString(si.OS.Architecture)
	sb.WriteString("]\n")

	// kernel
	sb.WriteString("kernel: ")
	sb.WriteString(si.Kernel.Release)
	sb.WriteString("\n")

	// CPU
	sb.WriteString("cpu: ")
	sb.WriteString(si.CPU.Model)
	sb.WriteString(" [")
	sb.WriteString(fmt.Sprint(si.CPU.Cores))
	sb.WriteString("/")
	sb.WriteString(fmt.Sprint(si.CPU.Threads))
	sb.WriteString("]\n")

	// Memory
	memory, err := memory.Get()
	sb.WriteString("memory: ")
	if err != nil {
		sb.WriteString("unknown")
	} else {
		sb.WriteString(fmt.Sprintf("%.2fGi used, %.2fGi free (%.2fGi total)", float64(memory.Used)/(1024*1024*1024), float64(memory.Free)/(1024*1024*1024), float64(memory.Total)/(1024*1024*1024)))
	}
	sb.WriteString("\n")

	// Storage
	sb.WriteString("storage:\n")
	for _, disk := range si.Storage {
		sb.WriteString("  ")
		sb.WriteString(disk.Name)
		sb.WriteString(" [")
		sb.WriteString(fmt.Sprint(disk.Model))
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprint(disk.Size))
		sb.WriteString("GB]\n")
	}

	// Network
	sb.WriteString("network:\n")
	for i, net := range si.Network {
		sb.WriteString("  ")
		sb.WriteString(net.Name)
		sb.WriteString(" [")
		sb.WriteString(fmt.Sprint(net.MACAddress))
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprint(net.Speed))
		sb.WriteString("Mbps]")
		if i < len(si.Network)-1 {
			sb.WriteString("\n")
		}
	}

	return WidgetResponse{
		"sysinfo",
		sb.String(),
		"",
	}, nil
}
