// ticket_handler.go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "strconv"
    "sync"

    "github.com/google/uuid"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    pb "your-github-repo/train-ticket-app/server"
)

type server struct {
    mu            sync.Mutex
    seatAllocations map[string]pb.SeatAllocation
}

func (s *server) PurchaseTicket(ctx context.Context, req *pb.TicketRequest) (*pb.TicketResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    receipt := generateReceipt()
    section := randomSection()
    seatNumber := randomSeatNumber()

    s.seatAllocations[req.UserEmail] = pb.SeatAllocation{
        UserEmail:   req.UserEmail,
        Section:     section,
        SeatNumber:  seatNumber,
    }

    return &pb.TicketResponse{Receipt: receipt}, nil
}

func (s *server) GetReceipt(ctx context.Context, req *pb.ReceiptRequest) (*pb.ReceiptResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    seatAllocation, exists := s.seatAllocations[req.UserEmail]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }

    return &pb.ReceiptResponse{Receipt: generateReceipt(seatAllocation)}, nil
}

func (s *server) GetSeatAllocation(ctx context.Context, req *pb.SeatRequest) (*pb.SeatResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    var seatAllocations []pb.SeatAllocation
    for _, allocation := range s.seatAllocations {
        if allocation.Section == req.Section {
            seatAllocations = append(seatAllocations, allocation)
        }
    }

    return &pb.SeatResponse{SeatAllocation: seatAllocations}, nil
}

func (s *server) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    _, exists := s.seatAllocations[req.UserEmail]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }

    delete(s.seatAllocations, req.UserEmail)
    return &pb.RemoveUserResponse{Message: "User removed successfully"}, nil
}

func (s *server) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    allocation, exists := s.seatAllocations[req.UserEmail]
    if !exists {
        return nil, fmt.Errorf("user not found")
    }

    allocation.Section = req.NewSection
    allocation.SeatNumber = req.NewSeatNumber

    s.seatAllocations[req.UserEmail] = allocation

    return &pb.ModifySeatResponse{Message: "Seat modified successfully"}, nil
}

func generateReceipt(seatAllocation ...pb.SeatAllocation) string {
    receiptID := uuid.New().String()
    var allocationText string

    if len(seatAllocation) > 0 {
        allocationText = fmt.Sprintf("Section: %s, Seat: %s", seatAllocation[0].Section, seatAllocation[0].SeatNumber)
    }

    return fmt.Sprintf("Receipt ID: %s\n%s", receiptID, allocationText)
}

func randomSection() string {
    sections := []string{"A", "B"}
    return sections[rand.Intn(len(sections))]
}

func randomSeatNumber() string {
    return strconv.Itoa(rand.Intn(10) + 1)
}

func main() {
    // Initialize gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterTicketServiceServer(s, &server{})

    // Register reflection service on gRPC server.
    reflection.Register(s)

    log.Printf("Server listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
