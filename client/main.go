// client/main.go
package main

import (
    "context"
    "fmt"
    "log"

    "google.golang.org/grpc"
    pb "your-github-repo/train-ticket-app/server"
)

func main() {
    // Set up gRPC connection
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil
