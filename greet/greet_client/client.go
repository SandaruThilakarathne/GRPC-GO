package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"

	"GRPC_GO_CLIENT/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm client")
	tls := false
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt"
		creds, sslErr:= credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Cannot load the client crt files: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50052", opts)
	if err != nil {
		log.Fatalf("couldn not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//println("Created: %v", c)
	//doUnary(c)
	//doCalculate(c)
	//doServerStreaming(c)
	//doPrimeStreaming(c)
	//doLongGreeting(c)
	//doAveraging(c)
	doBiDi(c)
	//doMaxNumber(c)
	//doErrorHandledSquarootUnary(c)
	//doUnaryCallDeadLine(c, 5*time.Second) // should complete
	//doUnaryCallDeadLine(c, 1*time.Second) // should timeout
}

func doUnary(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Theesh",
			LastName:  "Thilakarathne",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling greet RPC: %v", err)
	}

	log.Printf("Response from greet: %v", res.Result)

}

func doCalculate(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetSumRequest{
		Greetsum: &greetpb.GreetSum{
			FirstNumber: 3,
			LastNumber:  10,
		},
	}
	res, err := c.GreetCalculator(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling greet RPC: %v", err)
	}

	log.Printf("Response from greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Strting Server Streaming")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Theesh",
			LastName:  "Thilakarathne",
		},
	}
	resStram, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error occured %v", err)
	}

	for {
		msg, err := resStram.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stram %v", err)
		}
		log.Printf("Response from greet many times: %v", msg.GetResult())
	}
}

func doPrimeStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Strting Prime Streaming")
	req := &greetpb.PrimeNumberDecompositionRequest{
		Number: 100,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error occured %v", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stram %v", err)
		}

		log.Printf("Response from greet many times: %v", res.GetPrimeFactor())
	}
}

func doLongGreeting(c greetpb.GreetServiceClient) {
	fmt.Printf("Strting LongGreet Streaming")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Theesh",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Saja",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Paka",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while LongGreet Reading reading stram %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciveing response %v", err)
	}

	fmt.Printf("Longgreet resps: %v\n", res)

}

func doAveraging(c greetpb.GreetServiceClient) {
	fmt.Printf("Strting Averaging Streaming")

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error while LongGreet Reading reading stram %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}
	for _, number := range numbers {
		fmt.Printf("Sending request: %v\n", number)
		stream.Send(&greetpb.AverageRequest{
			Number: number,
		})
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciveing response %v", err)
	}

	fmt.Printf("Average: %v\n", res.GetAverage())
}

func doBiDi(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do BiDi Streaming RPC")
	// create stream by invoking client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Theesh",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Saja",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Paka",
			},
		},
	}

	waitc := make(chan struct{})
	// send a bunch of messages to the client (go routine)
	go func() {
		// function to send bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message:%v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive bunch of message from the server (go routing)
	go func() {
		// receive bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while recieving: %v", err)
				break
			}

			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doMaxNumber(c greetpb.GreetServiceClient) {
	fmt.Println("Max Number BiDi Streaming RPC")
	stream, err := c.MaxNumber(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	numbers := []int32{3, 5, 2, 9, 54, 23}

	waitc := make(chan struct{})

	go func() {
		for _,number := range numbers {
			//fmt.Printf("Sending message:%v\n", number)
			stream.Send(&greetpb.MaxNumberRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while recieving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorHandledSquarootUnary(c greetpb.GreetServiceClient)  {
	fmt.Printf("Starting squar root unary")
	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -10)

}

func doErrorCall(c greetpb.GreetServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &greetpb.SquareRootRequest{Number: n})
	if err != nil {
		responseErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC
			fmt.Println(responseErr.Message())
			if responseErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number\n")
			}
		} else {
			log.Fatalf("Big error calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("Results of square root number %v: %v\n", n, res.GetNumberRoot())
}

func doUnaryCallDeadLine(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Printf("Starting Unary with deadline")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Theesh",
			LastName:  "Thilakarathne",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Printf("Timeout was hit! Deadline was exeeded\n")
			}
		} else {
			log.Printf("Unexpecetd error: %v\n", statusErr)
		}
		return
	}

	log.Printf("Response from greet: %v", res.Result)
}