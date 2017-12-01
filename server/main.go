package main


import (
	"log"
	"net"

	pb "stschn.microsoft.com/logstreamer/pb"
	"github.com/hpcloud/tail"
	//"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) GetLogs(req *pb.LogPathRequest, stream pb.Logs_GetLogsServer) error {
	t, err := tail.TailFile(req.Path, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		stream.Send(&pb.LogEntryReply{Log: line.Text})
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLogsServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}