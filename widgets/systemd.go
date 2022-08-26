package widgets

import (
	"fmt"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/spf13/viper"
)

func SystemdWidget(v *viper.Viper, f formatFn) (WidgetResponse, error) {
	sb := strings.Builder{}

	f1 := f("7", "", false)
	f2 := f("10", "", true)
	f3 := f("1", "", true)
	f4 := f("11", "", false)

	conn, _ := dbus.New()
	defer conn.Close()

	units, _ := conn.ListUnitsFiltered([]string{"failed"})

	if len(units) > 0 {
		sb.WriteString(f3(fmt.Sprintf("%d", len(units))) + f1(" service(s) ") + f3("FAILED") + "\n")
		for i, unit := range units {
			sb.WriteString(f1("- ") + f4(fmt.Sprintf("%s", unit.Name)))
			if i < len(units)-1 {
				sb.WriteString("\n")
			}
		}
	} else {
		sb.WriteString(f1("All services ") + f2("OK"))
	}

	fmt.Println(units)

	return WidgetResponse{
		"systemd",
		sb.String(),
		"",
	}, nil
}
