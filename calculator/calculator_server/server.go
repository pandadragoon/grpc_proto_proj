package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pandadragoon/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Sum function was invoked with %v\n", req)
	first_number := req.GetFirstNumber()
	second_number := req.GetSecondNumber()

	sum_result := first_number + second_number

	res := &calculatorpb.SumResponse{
		SumResult: sum_result,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)
	number := req.GetNumber()
	var divisor int64 = 2

	for number > 1 {
		if number%divisor == 0 {
			err := stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: number,
			})
			if err != nil {
				return err
			}
			number = number / divisor
		} else {
			divisor++

			go func() {
				fmt.Printf("Divisor has been increased to: %v\n", divisor)
			}()
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	log.Println("ComputeAverage function was invoked")
	sum := float64(0)
	count := float64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			err := stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: sum / count,
			})
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to connect to client stream %v", err)
			return err
		}
		number := float64(req.GetNumber())
		count++
		sum += number
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	maximum := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error recieving client stream %v", err)
			return err
		}

		if maximum < req.GetNumber() {
			maximum = req.GetNumber()
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: maximum,
			})
			if err != nil {
				log.Fatalf("Error sending response to client %v", err)
				return err
			}
		}
	}

	return nil
}

func main() {
	fmt.Println("Calculator Server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Unable to serve: %v", err)
	}
}
