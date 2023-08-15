package monitor

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gen2brain/beeep"
	"google.golang.org/grpc/status"
)

func GetExternalIp() []byte {
	requestURL := "https://ipinfo.io/ip"
	client := http.Client{Timeout: 3 * time.Second}
	res, err := client.Get(requestURL)
	if err != nil {
		errcode, _ := status.FromError(err)
		errmsg := fmt.Sprintf("Unable to check external IP: %s", errcode.Message())
		beeep.Alert("NordLayer Helper", errmsg, "")
		log.Fatalf("[ExternalIP] client: could not connect to endpoint: %s ", err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[ExternalIP] client: could not read response body: %s\n", err)
	}
	log.Printf("[ExternalIP] Got response status: %d", res.StatusCode)
	log.Printf("[ExternalIP] Got response body: %d", resBody)
	return resBody
}
