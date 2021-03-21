package main

import (
	"context"
	"fmt"
	"github.com/pandadragoon/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello I'm a Calculator Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connnect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 4,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")
	var number int64 = 2342424534
	req := calculatorpb.PrimeNumberDecompositionRequest{
		Number: number,
	}

	client, err := c.PrimeNumberDecomposition(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling PrimeNumberDecomposition server %v", err)
	}

	for {
		res, err := client.Recv()
		if err == io.EOF {
			fmt.Printf("Done decomposing")
			break
		}
		if err != nil {
			log.Fatalf("Error getting sever stream %v", err)
		}

		fmt.Printf("Prime factor %d of %d\n", res.GetPrimeFactor(), number)
	}

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	numbers := []int32{
		20,
		30,
		22,
		26,
		37,
		18,
		25,
		25,
		26,
		34,
		21,
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error getting stream from server %v", err)
	}

	for _, num := range numbers {
		req := calculatorpb.ComputeAverageRequest{
			Number: num,
		}

		fmt.Printf("Sending requests %d \n", req.GetNumber())
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("Error sending request on stream %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error closing stream %v", err)
	}

	fmt.Printf("Average of %v was %f \n", numbers, res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error connections to server stream %v", err)
	}
	maximum := int32(0)

	waitc := make(chan struct{})

	numbers := []int32{
		20,
		30,
		22,
		26,
		37,
		18,
		25,
		25,
		26,
		34,
		21,
	}

	go func() {
		for _, num := range numbers {
			req := calculatorpb.FindMaximumRequest{
				Number: num,
			}

			err := stream.Send(&req)
			if err == io.EOF {
				log.Printf("Done streaming to server stream")
				break
			}
			if err != nil {
				log.Fatalf("Error sending request to server stream %v", err)
			}

			time.Sleep(500 * time.Millisecond)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Error closing stream %v", err)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Err receiving server stream %v", err)
			}

			maximum = res.GetMaximum()
			log.Printf("received maximum %d", maximum)
		}

		close(waitc)
	}()

	<-waitc
}
