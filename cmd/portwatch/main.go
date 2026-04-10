package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func main() {
	cfg, err := config.Resolve(os.Args[1:])
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	sc, err := scanner.NewScanner(cfg.StartPort, cfg.EndPort)
	if err != nil {
		log.Fatalf("scanner error: %v", err)
	}

	logHandler := alert.NewLogHandler(log.New(os.Stdout, "", log.LstdFlags), "")
	dispatcher := alert.NewDispatcher()
	dispatcher.Register(logHandler)

	mon := monitor.New(sc, dispatcher)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("shutting down portwatch...")
		cancel()
	}()

	interval := time.Duration(cfg.IntervalSecs) * time.Second
	log.Printf("portwatch started: scanning ports %d-%d every %s",
		cfg.StartPort, cfg.EndPort, interval)

	if err := mon.Run(ctx, interval); err != nil {
		log.Fatalf("monitor error: %v", err)
	}
}
