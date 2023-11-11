package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"github.com/ihyoudou/nordlayer-helper/api"
	"github.com/ihyoudou/nordlayer-helper/icon"
	"google.golang.org/grpc/status"
)

const (
	APP_NAME    = "Nordlayer Helper"
	CHECK_EVERY = 5 * time.Second
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Red, icon.Red)
	systray.SetTitle(APP_NAME)

	// inital menu items
	ConnectedmenuItem := systray.AddMenuItem("Connected: No", "Logged in?")
	InternalIPmenuItem := systray.AddMenuItem(`Internal IP: ?`, "Internal network IP")
	ExternalIPmenuItem := systray.AddMenuItem(`Server IP: ?`, "IP of vpn server endpoint")
	ProtocolmenuItem := systray.AddMenuItem(`Protocol: ?`, "VPN Protocol")

	// quit function
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// Get available gateways
	res, err := api.GetVPNGateways()
	if err != nil {
		log.Printf("issue while asking nordlayer for available gateways: %s", err)
		errcode, _ := status.FromError(err)

		if errcode.Message() == "user_id not set" {
			systray.SetTemplateIcon(icon.Red, icon.Red)
			beeep.Alert(APP_NAME, "user_id not set - you got logged out", "")
		}
	}

	for index, element := range res.AvailableGateways {
		ConnectedmenuItem.AddSubMenuItem(element, fmt.Sprintf("Connect to %s (id: %s)", element, index))
	}

	// Refreshing nordlayer status
	go func() {
		tick := time.Tick(CHECK_EVERY)
		unableToConnectAlert := false
//		lastIpCheck := "0.0.0.0"
		for range tick {

			status, err := api.GetVPNStatus()
			if err != nil {
				log.Printf("issue while asking nordlayer for status: %s", err)
				if unableToConnectAlert == false {
					beeep.Alert(APP_NAME, "Was unable to connect to NordLayer Backend - make sure it is running", "")
					unableToConnectAlert = true
				}
			} else {
				if unableToConnectAlert == true {
					beeep.Alert(APP_NAME, "Was enable to connect back to NordLayer Backend", "")
				}
				// Clearing error states
				unableToConnectAlert = false
			}

			// Checking external IP
//			externalIp := monitor.GetExternalIp()
//			log.Printf("lastip: %s", lastIpCheck)
//			log.Printf("extip: %s", externalIp)
//			if lastIpCheck != externalIp {
//				if lastIpCheck != "0.0.0.0" {
//					msg := fmt.Sprintf("ExternalIP has changed from %s to %s", lastIpCheck, externalIp)
//					log.Printf(msg)
//					beeep.Alert(APP_NAME, msg, "")
//				}
//				lastIpCheck = externalIp
//			}

			// if connected gateway is empty (eg we are not connected), setting placeholder values
			if status.ConnectedGateway == "" {
				systray.SetTemplateIcon(icon.Red, icon.Red)
				ConnectedmenuItem.SetTitle("Connected: No")
				InternalIPmenuItem.SetTitle(`Internal IP: ?`)
				ProtocolmenuItem.SetTitle(`Protocol: ?`)
			} else {
				systray.SetTemplateIcon(icon.Green, icon.Green)
				ConnectedmenuItem.SetTitle(fmt.Sprintf("Connected: %s", status.ConnectedGateway))
				InternalIPmenuItem.SetTitle(fmt.Sprintf("Internal IP: %s", status.InternalIp))
				ProtocolmenuItem.SetTitle(fmt.Sprintf("Protocol: %s", status.Protocol))
			}

			ExternalIPmenuItem.SetTitle(fmt.Sprintf("Server IP: %s", status.ServerIp))

		}
	}()

}

func onExit() {}
