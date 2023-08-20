package api

import (
	"context"
	"log"
	"time"

	pb "github.com/ihyoudou/nordlayer-helper/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	ADDRESS = "unix:///run/nordlayer/nordlayer.sock"
	TIMEOUT = time.Second
)

func GetVPNStatus() (*pb.Payload, error) {

	log.Print("[vpnstatus] Trying to connect")
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	conn, err := grpc.DialContext(ctx, ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Printf("did not connect: %s", err)
		return &pb.Payload{}, err
	} else {
		defer conn.Close()
	}

	log.Print("Trying to create new daemon client")
	c := pb.NewStatusClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	log.Print("[vpnstatus] Trying to get status")
	res, err := c.Get(ctx, &pb.Empty{})

	if err != nil {
		log.Fatalf("could not send: %v", err)
		code := status.Code(err)
		log.Fatalf("got error: %v", code)
	}
	log.Print(res.Payload)
	return res.Payload, err
}

func GetVPNGateways() (*pb.Gateways, error) {

	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	conn, err := grpc.DialContext(ctx, ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Printf("did not connect: %s", err)
		return &pb.Gateways{}, err
	} else {
		defer conn.Close()
	}
	log.Print("Trying to create new daemon client")
	c := pb.NewVpnClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	log.Print("[vpngateways] Trying to get status")
	res, err := c.Gateways(ctx, &pb.Empty{})

	if err != nil {
		code, ok := status.FromError(err)
		log.Print(ok)
		log.Fatalf("got error: %v", code.Message())
		log.Fatalf("vpngateways: could not send: %v", err)
		// user_id not set
	}

	log.Print(res.Gateways)
	return res.Gateways, err
}

// func VPNConnect(gateway string) error {

// 	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
// 	conn, err := grpc.DialContext(ctx, ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

// 	if err != nil {
// 		log.Printf("did not connect: %s", err)
// 		return &pb.Payload{}, err
// 	} else {
// 		defer conn.Close()
// 	}

// 	log.Print("Trying to create new daemon client")
// 	c := pb.NewVpnClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

// 	defer cancel()
// 	log.Print("Trying to get status")
// 	res, err := c.Connect(ctx, &pb.ConnectRequest{GatewayName: gateway})

// 	if err != nil {
// 		log.Fatalf("could not send: %v", err)
// 	}

// 	log.Print(res)
// 	return err
// }
