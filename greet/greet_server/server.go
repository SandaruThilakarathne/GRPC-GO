package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"GRPC_GO/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)
	firstName := req.Greeting.GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

// sum
func (*server) GreetCalculator(ctx context.Context, req *greetpb.GreetSumRequest) (*greetpb.GreetSumResponse, error) {
	fmt.Printf("GreetSum function was invoked with %v", req)
	firstNumber := req.Greetsum.GetFirstNumber()
	lastNumber := req.Greetsum.GetLastNumber()

	result := firstNumber + lastNumber

	res := &greetpb.GreetSumResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet many times function was invoked %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number" + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (*server) PrimeNumberDecomposition(req *greetpb.PrimeNumberDecompositionRequest, stream greetpb.GreetService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked %v\n", req)

	number := req.GetNumber()

	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&greetpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor += 1
			fmt.Printf("Divisor incremented to %v\n\n", divisor)
		}
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function streaming %v\n", stream)
	result := "Hello"

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// reading finished
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})

		}
		if err != nil {
			log.Fatalf("Stucked with an error: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "

	}
}

func (*server) Average(stream greetpb.GreetService_AverageServer) error {
	fmt.Printf("Average function streaming %v\n", stream)
	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// reading finished
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&greetpb.AverageResponse{
				Average: average,
			})
		}

		if err != nil {
			log.Fatalf("Stucked with an error: %v", err)
			return err
		}

		sum += req.GetNumber()
		count++
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("Greet everyone function streaming %v\n", stream)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// reading finished
			return nil
		}
		if err != nil {
			log.Fatalf("Stucked with an error: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello" + firstName + "!"
		sendErr := stream.Send(&greetpb.GreetEvryonResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}

	}
}

func (*server) MaxNumber(stream greetpb.GreetService_MaxNumberServer) error {
	fmt.Printf("MaxNumber function streaming %v\n", stream)
	maximum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Stucked with an error: %v", err)
		}
		number := req.GetNumber()
		if number > maximum {
			maximum = number

			sendErr := stream.Send(&greetpb.MaxNumberResponse{
				Result: number,
			})

			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *greetpb.SquareRootRequest) (*greetpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number))

	}

	return &greetpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDealine function was invoked with %v", req)
	for i:=0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client cancelled the request")
			return nil, status.Error(codes.DeadlineExceeded, "the client cancelled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.Greeting.GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Hello world")
	opts := []grpc.ServerOption{}
	lis, err := net.Listen("tcp", "0.0.0.0:3000")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	tls := false
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed to loading certificates: %v", sslErr)
		}
		opts = append(opts, grpc.Creds(creds))
	}
	s := grpc.NewServer(opts...)
	reflection.Register(s)
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
