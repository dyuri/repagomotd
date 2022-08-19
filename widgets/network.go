package widgets

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/viper"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func NetworkWidget(v *viper.Viper, f formatFn) (WidgetResponse, error) {
	sb := strings.Builder{}

	f1 := f("7", "", false)
	f2 := f("11", "", false)
	f3 := f("2", "", false)
	f4 := f("10", "", true)

	addLine := func(title, value string) {
		fmt.Fprintf(&sb, "%s %s\n", f1(fmt.Sprintf("%9s", title)), value)
	}

	excluded := v.GetStringSlice("net.exclude")
	included := v.GetStringSlice("net.include")

	netInterfaces, _ := net.Interfaces()
	netInterfaceStats, _ := net.IOCounters(true)
	for _, netInterface := range netInterfaces {
		ifname := netInterface.Name
		if (len(included) > 0 && !stringInSlice(ifname, included)) || (len(excluded) > 0 && stringInSlice(ifname, excluded)) {
			continue
		}
		sent := 0
		recv := 0
		errors := 0
		for _, netInterfaceStat := range netInterfaceStats {
			if netInterfaceStat.Name == ifname {
				sent = int(netInterfaceStat.BytesSent / 1024 / 1024)
				recv = int(netInterfaceStat.BytesRecv / 1024 / 1024)
				errors = int(netInterfaceStat.Errout + netInterfaceStat.Errin)
			}
		}
		stats := ""

		// TODO ipv6 based on config?
		for _, addr := range netInterface.Addrs {
			if strings.Contains(addr.Addr, ".") {
				stats += f4(fmt.Sprintf("%18s", addr.Addr)) + " "
			}
		}

		stats += f3(fmt.Sprintf("%dM/%dM ⭿", recv, sent))
		if errors > 0 {
			stats += f2(fmt.Sprintf(" %d ⚠", errors))
		}
		addLine(netInterface.Name, stats)
	}

	return WidgetResponse{
		"net",
		sb.String(),
		"",
	}, nil
}
