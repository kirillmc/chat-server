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
var examplePath = flag.String("ex", "", "example path")

const (
	servicePort   = 50051
	ExamplePath   = "/user_v1.UserV1/Get"
	SERVICE_PEM   = "tls/client/service.pem"
	AUTHORIZATION = "Authorization"
	REARER        = "Bearer "
)

func main() {
	flag.Parse()

	ctx := context.Background()
	md := metadata.New(map[string]string{AUTHORIZATION: REARER + *accessToken})
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
		EndpointAddress: *examplePath,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Access granted")
}

/** For tests:
go run cmd/grpc_client/main.go -a "" -ex "/user_v1.UserV1/Get"
*/
