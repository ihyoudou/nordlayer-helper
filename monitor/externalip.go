package monitor

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"google.golang.org/grpc/status"
)

var lastUsedService = 0

func GetExternalIp() string {
	// Using array of multiple services that returns external IP to not hammer one
	externalIpServices := []string{"https://ipinfo.io/ip", "https://ifconfig.me", "https://icanhazip.com", "http://ipecho.net/plain", "https://ident.me", "https://api.ipify.org"}
	// Going back to first element
	if lastUsedService >= len(externalIpServices) {
		lastUsedService = 0
	}

	requestURL := externalIpServices[lastUsedService]
	log.Printf("[ExternalIP] Using service: %s", externalIpServices[lastUsedService])

	client := http.Client{Timeout: 2 * time.Second}
	res, err := client.Get(requestURL)
	if err != nil {
		errcode, _ := status.FromError(err)
		errmsg := fmt.Sprintf("Unable to check external IP: %s", errcode.Message())
		beeep.Alert("NordLayer Helper", errmsg, "")
		log.Printf(errmsg)
		return "0.0.0.0"
	} else {
		defer res.Body.Close()
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[ExternalIP] client: could not read response body: %s\n", err)
	}
	log.Printf("[ExternalIP] Got response status: %d", res.StatusCode)
	log.Printf("[ExternalIP] Got response body: %s", string(resBody))

	lastUsedService++
	return strings.TrimSpace(string(resBody))
}
