package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	pb "server/proto"
	"time"
)

type server struct {
	pb.UnimplementedSappServiceServer
}

// Handles incoming batch requests
func (s *server) SmSPacket(ctx context.Context, in *pb.SmSRequest) (*pb.SmSResponse, error) {
	log.Printf("Received SmSRequest with %d messages", len(in.Messages))

	// Process each incoming message
	var responseMessages []*pb.SmS
	for _, message := range in.Messages {
		log.Printf("Processing message: From=%v, To=%v, Amount=%v, OperationTime=%v",
			message.From, message.To, message.Amount, message.OperationTime)

		// For this example, we'll add a suffix to the OperationTime to indicate it's been processed
		responseMessages = append(responseMessages, &pb.SmS{
			From:          message.From,
			To:            message.To,
			Amount:        message.Amount,
			OperationTime: time.Now().Format(time.RFC3339) + " - processed",
		})
	}

	// Return the response containing all processed messages
	return &pb.SmSResponse{Messages: responseMessages}, nil
}

type Configs struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	Protocol string `yaml:"protocol"`
}

var (
	conf     = Configs{}
	clientID string
	InfoLog  *log.Logger
)

func config() {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}
	fmt.Println("Raw YAML content:", string(data))
	// Unmarshal the YAML data into a Config struct
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
	if conf.Port == "" || conf.Host == "" {
		log.Println("One or more fields are not populated.")
	}
}
func init() {
	config()
	logData()
}
func main() {

	fmt.Println("gRPC server is running on port ", conf.Port, "\n ")

	listen()

}

func listen() {
	listener, err := net.Listen(conf.Protocol, conf.Host+conf.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//grpcServer
	gs := grpc.NewServer()
	pb.RegisterSappServiceServer(gs, &server{})

	if err := gs.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
func logData() {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLog = log.New(logFile, "Info "+time.Now().Format(time.RFC3339Nano)+clientID+"\n", log.Ldate|log.Ltime|log.Lshortfile)
}
