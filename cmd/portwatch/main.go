package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.github.com/user/portwatch"
	"github./scanner"
	"github.com/user/portwatch/internal/state"
)

func main() {
	cfg, err := config.Resolve(os.Args[1:])
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	sc, err := scanner.NewScanner(cfg.PortRange)
	if err != nil {
		log.Fatalf("scanner: %v", err)
	}

	store := state.NewStore(cfg.StateFile)
	mgr := state.NewManager(store)
	printer := report.NewPrinter(os.Stdout)

	var handlers []alert.Handler
	handlers = append(handlers, alert.NewLogHandler("[portwatch] "))
	if cfg.Webhook != "" {
		handlers = append(handlers, alert.NewWebhookHandler(cfg.Webhook))
	}
	if cfg.Email.Enabled {
		emailCfg := alert.EmailConfig{
			Host:     cfg.Email.Host,
			Port:     cfg.Email.Port,
			Username: cfg.Email.Username,
			Password: cfg.Email.Password,
			From:     cfg.Email.From,
			To:       cfg.Email.To,
		}
		handlers = append(handlers, alert.NewEmailHandler(emailCfg))
	}

	ch := make(chan []monitor.Change, 8)
	dispatcher := alert.NewDispatcher(ch, alert.NewMultiHandler(handlers...))
	go dispatcher.Run()

	ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
	defer ticker.Stop()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		ports, err := sc.Scan()
		if err != nil {
			log.Printf("scan error: %v", err)
		} else {
			changes, err := mgr.Update(ports)
			if err != nil {
				log.Printf("state update error: %v", err)
			} else {
				printer.PrintSnapshot(ports, changes)
				if len(changes) > 0 {
					ch <- changes
				}
			}
		}
		select {
		case <-ticker.C:
		case <-sig:
			log.Println("shutting down")
			return
		}
	}
}
