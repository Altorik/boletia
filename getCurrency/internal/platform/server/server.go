package server

import (
	"boletia/internal/obtain"
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"boletia/internal/platform/server/handler/currency"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpAddr        string
	engine          *gin.Engine
	getCurrency     obtain.CurrencyService
	shutdownTimeout time.Duration
}

//func Run(ctx context.Context, tickDuration time.Duration, getCurrency obtain.CurrencyService) error {
//	serverContext := serverContext(ctx)
//	defer panicHandler()
//	return currency.ObtainHandler(serverContext, getCurrency)
//}

func Run(ctx context.Context, tickDuration time.Duration, getCurrency obtain.CurrencyService) error {
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()
	serverContext := serverContext(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Println("Terminating periodic task")
			return nil
		case <-ticker.C:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic: %v", r)
					}
				}()

				err := currency.ObtainHandler(serverContext, getCurrency)
				if err != nil {
					log.Printf("Error fetching data: %v", err)
				} else {
					log.Println("Data fetched successfully")
				}
			}()
		}
	}
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
