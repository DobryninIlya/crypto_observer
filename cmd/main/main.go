package main

import (
	"context"
	application "cryptoObserver/internal/app/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := application.LoadConfig()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var app *application.App
	var err error
	go handleSignals(cancel)
	go func() {
		defer cancel()
		log.Printf("Запуск сервера на %s", config.Server.Port)
		log.Printf("Config: %+v", config)
		if app, err = application.Start(ctx, config); err != nil {
			log.Fatal(err)
		}

		if err = app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			if err = app.Server.Shutdown(ctx); err != nil {
				log.Printf("Ошибка при остановке сервера: %v", err)
				if err = app.Server.Close(); err != nil {
					log.Printf("Ошибка при закрытии сервера: %v", err)
				}
				return
			} else {
				if err = app.Server.Close(); err != nil {
					log.Printf("Ошибка при закрытии сервера: %v", err)
				}
				log.Println("Сервер успешно остановлен")
				return
			}
		}
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	cancel()
}
