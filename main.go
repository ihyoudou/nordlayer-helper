package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/ihyoudou/nordlayer-helper/api"
	"github.com/ihyoudou/nordlayer-helper/icon"
	"github.com/ihyoudou/nordlayer-helper/monitor"
	"google.golang.org/grpc/status"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Red, icon.Red)
	systray.SetTitle("Nordlayer Helper")

	// inital menu items
	ConnectedmenuItem := systray.AddMenuItem("Connected: No", "Logged in?")
	InternalIPmenuItem := systray.AddMenuItem(`IntIP: ?`, "Internal IP")
	ExternalIPmenuItem := systray.AddMenuItem(`ExtIP: ?`, "External IP")
	ProtocolmenuItem := systray.AddMenuItem(`Protocol: ?`, "VPN Protocol")

	lastAlertSend := ""
	// Get available gateways
	res, err := api.GetVPNGateways()
	if err != nil {
		log.Printf("issue while asking nordlayer for available gateways: %s", err)
		errcode, ok := status.FromError(err)
		log.Print(ok)
		if (errcode.Message() == "user_id not set") && (lastAlertSend != "user_id not set") {
			systray.SetTemplateIcon(icon.Red, icon.Red)
			beeep.Alert("NordLayer Helper", "user_id not set - you got logged out", "")
			lastAlertSend = "user_id not set"
		}
	} else {
		lastAlertSend = ""
	}

	for index, element := range res.AvailableGateways {
		ConnectedmenuItem.AddSubMenuItem(element, fmt.Sprintf("Connect to %s (id: %s)", element, index))
	}

	// Refreshing nordlayer status
	go func() {

		tick := time.Tick(2000 * time.Millisecond)
		for range tick {

			status, err := api.GetVPNStatus()
			if err != nil {
				log.Printf("issue while asking nordlayer for status: %s", err)
			}
			// if connected gateway is empty (eg we are not connected), setting placeholder values
			if status.ConnectedGateway == "" {
				externalIp := monitor.GetExternalIp()
				pingStats := monitor.PingHost("8.8.8.8")
				log.Print(pingStats.AvgRtt)
				systray.SetTemplateIcon(icon.Red, icon.Red)
				ConnectedmenuItem.SetTitle("Connected: No")
				InternalIPmenuItem.SetTitle(`IntIP: ?`)
				ExternalIPmenuItem.SetTitle(fmt.Sprintf("ExtIP: %s", externalIp))
				ProtocolmenuItem.SetTitle(`Protocol: ?`)
			} else {
				systray.SetTemplateIcon(icon.Green, icon.Green)
				ConnectedmenuItem.SetTitle(fmt.Sprintf("Connected: %s", status.ConnectedGateway))
				InternalIPmenuItem.SetTitle(fmt.Sprintf("IntIP: %s", status.InternalIp))
				ExternalIPmenuItem.SetTitle(fmt.Sprintf("ExtIP: %s", status.ExternalIp))
				ProtocolmenuItem.SetTitle(fmt.Sprintf("Protocol: %s", status.Protocol))
			}

			// pingStats := monitor.pingHost("8.8.8.8")

		}
	}()

	// quit function
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

}

func onExit() {}
