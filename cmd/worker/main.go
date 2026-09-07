package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"payment_gateway/internal/config"
	"payment_gateway/internal/kafka/consumer"
	applogger "payment_gateway/internal/lib/logger"
	"payment_gateway/internal/payment"
	"payment_gateway/internal/storage/postgres"
	"payment_gateway/internal/worker"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error(
			"worker stopped",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env: %v", err)
	}

	cfg := config.LoadConfig()
	logger := applogger.SettupLogger(cfg.Env)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	storage, err := postgres.SetupDatabase(
		ctx,
		cfg.Storage.URL(),
		10,
	)
	if err != nil {
		return fmt.Errorf("setup database: %w", err)
	}
	defer storage.Close()

	paymentRepository := postgres.NewPaymentRepository(storage)
	paymentService := payment.NewService(paymentRepository)

	paymentCreateHandler :=
		consumer.NewPaymentCreateHandler(paymentService)

	paymentStatusHandler :=
		consumer.NewPaymentStatusHandler(paymentService)

	dispatcher := worker.NewDispatcher()

	if err := dispatcher.Register(
		cfg.Kafka.PaymentCreateTopic,
		paymentCreateHandler,
	); err != nil {
		return fmt.Errorf(
			"register payment create handler: %w",
			err,
		)
	}

	if err := dispatcher.Register(
		cfg.Kafka.PaymentStatusTopic,
		paymentStatusHandler,
	); err != nil {
		return fmt.Errorf(
			"register payment status handler: %w",
			err,
		)
	}

	kafkaClient, err := consumer.NewClient(
		cfg.Kafka.Brokers,
		cfg.Kafka.ConsumerGroup,
		[]string{
			cfg.Kafka.PaymentCreateTopic,
			cfg.Kafka.PaymentStatusTopic,
		},
	)
	if err != nil {
		return fmt.Errorf("create Kafka consumer: %w", err)
	}
	defer kafkaClient.Close()

	backgroundWorker := worker.New(
		kafkaClient,
		dispatcher,
	)

	logger.Info(
		"starting worker",
		slog.String("consumer_group", cfg.Kafka.ConsumerGroup),
	)

	if err := backgroundWorker.Run(ctx); err != nil {
		return fmt.Errorf("run worker: %w", err)
	}

	logger.Info("worker stopped")
	return nil
}
