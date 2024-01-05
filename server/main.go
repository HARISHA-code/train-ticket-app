package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "your-github-repo/train-ticket-app/server"
)

type server struct {
	mu              sync.Mutex
	seatAllocations map[string]pb.SeatAllocation
}
