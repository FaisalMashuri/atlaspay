package main

import (
	"auth_service/cmd"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := cmd.Run(ctx); err != nil {
		log.Printf("application stopped with error: %v", err)
		os.Exit(1)
	}

	log.Println("application stopped gracefully")
}
