package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/tima-fey/anti_brute/internal/localDB"
	"github.com/tima-fey/anti_brute/internal/scheme"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type db localDB.BaseDB // can't implement next line without this. why?

func (db db) CheckAll(ctx context.Context, r *scheme.Request) (*scheme.Answer, error) {
	log.Print(fmt.Sprintf("Got request CheckAll %v", r))
	answers := make(chan bool, 3)
	go db.Address.Add(r.Address, answers)
	go db.Login.Add(r.Login, answers)
	go db.Password.Add(r.Password, answers)
	answer1, answer2, answer3 := <-answers, <-answers, <-answers
	if answer1 && answer2 && answer3 {
		return &scheme.Answer{Allow: true}, nil
	}
	return &scheme.Answer{Allow: false}, nil
}

func StartGrpcServer(db db, address string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%d", address, port))
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	scheme.RegisterEventsHandlerServer(grpcServer, db)
	grpcServer.Serve(lis)
}
