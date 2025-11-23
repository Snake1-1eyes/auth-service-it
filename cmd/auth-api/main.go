package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Snake1-1eyes/auth-service-it/internal/config"
	delivery "github.com/Snake1-1eyes/auth-service-it/internal/delivery/grpc"
	"github.com/Snake1-1eyes/auth-service-it/internal/logger"
	"github.com/Snake1-1eyes/auth-service-it/internal/repository/postgres"
	"github.com/Snake1-1eyes/auth-service-it/internal/repository/redis"
	"github.com/Snake1-1eyes/auth-service-it/internal/usecase/auth"
	pb "github.com/Snake1-1eyes/auth-service-it/pkg/auth"
	_ "github.com/jackc/pgx/v5/stdlib"
	redisClient "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// Logger
	log, err := logger.New(cfg.Environment, cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	// Postgres
	db, err := sql.Open("pgx", cfg.GetPostgresDSN())
	if err != nil {
		// log.Fatal expects context, msg, fields
		// Since we don't have context yet, we can pass nil or context.Background()
		// But logger implementation checks for nil context?
		// Let's assume it handles nil or use context.Background()
		panic(fmt.Sprintf("failed to connect to postgres: %v", err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("failed to ping postgres: %v", err))
	}

	// Redis
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr: "redis:6379",
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	sessionRepo := redis.NewSessionRepository(rdb)

	// Usecase
	uc := auth.NewUseCase(userRepo, sessionRepo)

	// Delivery
	server := delivery.NewServer(uc)

	// gRPC Server
	lis, err := net.Listen("tcp", cfg.GetGRPCAddress())
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, server)
	reflection.Register(s)

	log.Info(context.Background(), fmt.Sprintf("server listening at %v", lis.Addr()))

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("failed to serve: %v", err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	s.GracefulStop()
	log.Info(context.Background(), "server stopped")
}
