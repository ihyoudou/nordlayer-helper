package monitor

import (
	"fmt"
	"log"

	"github.com/gen2brain/beeep"
	probing "github.com/prometheus-community/pro-bing"
	"google.golang.org/grpc/status"
)

func PingHost(host string) *probing.Statistics {

	pinger, err := probing.NewPinger(host)
	if err != nil {
		log.Fatalf("[Ping] could not start new pinger: %s ", err)
	}
	pinger.Count = 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		errcode, _ := status.FromError(err)
		errmsg := fmt.Sprintf("Unable to ping: %s", errcode.Message())
		beeep.Alert("NordLayer Helper", errmsg, "")
		log.Fatalf("[Ping] could not connect: %s ", err)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	return stats
}
