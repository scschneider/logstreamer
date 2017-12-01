package main


import (
	"os"
	"log"
	"io"
	"fmt"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "stschn.microsoft.com/logstreamer/pb"
)

const (
	port = ":50051"
	LP = "/var/log/system.log"
)

func GetLogs(client pb.LogsClient) {
	stream, err := client.GetLogs(context.Background(), &pb.LogPathRequest{Path: LP})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.Log)
	}
}

func main() {
	addresses := strings.Split(os.Args[1],",")

	for _, element := range addresses {
		conn, err := grpc.Dial(element + port, grpc.WithInsecure())
		if err != nil {
			fmt.Println(element)
			log.Fatal(err)
		}
		defer conn.Close()

		c := pb.NewLogsClient(conn)

		fmt.Println(element)
		GetLogs(c)
	}
}


