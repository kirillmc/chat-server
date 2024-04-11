package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	descAccess "github.com/kirillmc/auth/pkg/access_v1"
)

var accessToken = flag.String("a", "", "access token")

const (
	servicePort = 50051
	ExamplePath = "/user_v1.UserV1/Get"
	SERVICE_PEM = "tls/client/service.pem"
)

func main() {
	flag.Parse()

	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer " + *accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	creds, err := credentials.NewClientTLSFromFile(SERVICE_PEM, "")
	if err != nil {
		log.Fatalf("FAILED TO GET CREDS: %v\n", err)
	}

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", servicePort),
		grpc.WithTransportCredentials(creds),
		//grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial GRPC client: %v", err)
	}

	cl := descAccess.NewAccessV1Client(conn)

	_, err = cl.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: ExamplePath,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Access granted")
}
