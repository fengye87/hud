package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/libnetwork/resolvconf"
	"github.com/docker/libnetwork/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vishvananda/netlink"
	"gopkg.in/ini.v1"
)

const JournalMaxLines = 120

func main() {
	release, err := getOSRelease()
	if err != nil {
		panic(err)
	}

	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { return event })
	if err := setupSignalHandler(app); err != nil {
		panic(err)
	}

	headerTV := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorTeal).SetChangedFunc(func() { app.Draw() })
	headerTV.SetText("\n" + release)

	infoTV := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() { app.Draw() })

	journalTV := tview.NewTextView().ScrollToEnd().SetChangedFunc(func() { app.Draw() })
	journalTV.SetMaxLines(JournalMaxLines)
	journalctl := exec.Command("journalctl", "--no-hostname", "-n", strconv.Itoa(JournalMaxLines), "-f")
	journalctl.Stdout = journalTV
	if err := journalctl.Start(); err != nil {
		panic(err)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(headerTV, 3, 0, false).AddItem(infoTV, 0, 0, false).AddItem(journalTV, 0, 1, false)

	go func() {
		hostnameInfo := []string{"Hostname", ""}
		ipInfo := []string{"IP address", ""}
		gatewayInfo := []string{"Gateway address", ""}
		dnsInfo := []string{"DNS nameservers", ""}
		timeInfo := []string{"Time", ""}
		infoItems := [][]string{hostnameInfo, ipInfo, gatewayInfo, dnsInfo, timeInfo}
		lastUpdateTime := time.Time{}

		for {
			if time.Now().After(lastUpdateTime.Add(time.Minute)) {
				hostname, _ := getHostname()
				hostnameInfo[1] = hostname

				link, _ := getDefaultLink()
				if link != nil {
					addr, _ := getLinkAddr(link)
					if addr != nil {
						ipInfo[1] = addr.IPNet.String()
					}

					gateway, _ := getLinkGateway(link)
					if gateway != nil {
						gatewayInfo[1] = gateway.String()
					}
				}

				nameservers, _ := getNameservers()
				dnsInfo[1] = strings.Join(nameservers, " ")
			}

			timeInfo[1] = time.Now().Format(time.RFC1123Z)

			_, _, infoTVWidth, _ := infoTV.GetRect()
			var infoLines []string
			for _, info := range infoItems {
				label := info[0]
				paddingLen := infoTVWidth/2 - len(label) - 2
				if paddingLen < 0 {
					paddingLen = 0
				}

				value := info[1]
				if value == "" {
					value = "[:red]NOT SET[:-]"
				}

				infoLines = append(infoLines, fmt.Sprintf("%s[::b]%s[::-]: %s", strings.Repeat(" ", paddingLen), label, value))
			}

			infoTV.SetText(strings.Join(infoLines, "\n"))
			flex.ResizeItem(infoTV, len(infoLines)+1, 0)

			time.Sleep(time.Second)
		}
	}()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func setupSignalHandler(app *tview.Application) error {
	printkPath := "/proc/sys/kernel/printk"
	printk, err := os.ReadFile(printkPath)
	if err != nil {
		return fmt.Errorf("get printk: %s", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sigs

		os.WriteFile(printkPath, printk, 0644)
		app.Stop()
	}()

	if err := os.WriteFile(printkPath, []byte("0 0 0 0"), 0644); err != nil {
		return fmt.Errorf("set printk: %s", err)
	}
	return nil
}

func getOSRelease() (string, error) {
	info, err := ini.Load("/etc/os-release")
	if err != nil {
		return "", fmt.Errorf("read OS info: %s", err)
	}
	return info.Section("").Key("PRETTY_NAME").String(), nil
}

func getHostname() (string, error) {
	return os.Hostname()
}

func getDefaultLink() (netlink.Link, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, fmt.Errorf("list links: %s", err)
	}

	for _, link := range links {
		routes, err := netlink.RouteList(link, netlink.FAMILY_V4)
		if err != nil {
			return nil, fmt.Errorf("list routes of link %q: %s", link.Attrs().Name, err)
		}

		for _, route := range routes {
			if route.Dst == nil {
				return link, nil
			}
		}
	}
	return nil, nil
}

func getLinkAddr(link netlink.Link) (*netlink.Addr, error) {
	addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return nil, fmt.Errorf("list addrs: %s", err)
	}

	if len(addrs) > 0 {
		return &addrs[0], nil
	}
	return nil, nil
}

func getLinkGateway(link netlink.Link) (net.IP, error) {
	routes, err := netlink.RouteList(link, netlink.FAMILY_V4)
	if err != nil {
		return nil, fmt.Errorf("list routes: %s", err)
	}

	if len(routes) > 0 {
		return routes[0].Gw, nil
	}
	return nil, nil
}

func getNameservers() ([]string, error) {
	rc, err := resolvconf.Get()
	if err != nil {
		return nil, fmt.Errorf("get resolvconf: %s", err)
	}
	return resolvconf.GetNameservers(rc.Content, types.IPv4), nil
}
