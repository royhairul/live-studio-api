package jobs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/params"
	"github.com/royhairul/live-studio-api/internal/domains/transaction/service"
	"go.uber.org/fx"
)

type TransactionJobParams struct {
	fx.In
	TransactionSvc service.TransactionService
	Lifecycle      fx.Lifecycle
}

func StartTransactionJob(p TransactionJobParams) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0o755); err != nil {
		fmt.Fprintf(gin.DefaultWriter, "[CRON][ERROR] Failed to create logs directory: %v\n", err)
		return
	}

	// Open or create log file for cronjob
	logFile, err := os.OpenFile("logs/cronjob.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		fmt.Fprintf(gin.DefaultWriter, "[CRON][ERROR] Failed to open log file: %v\n", err)
		return
	}

	// Create a dedicated logger for cronjob
	cronLogger := log.New(logFile, "[TRANSACTION-JOB] ", log.LstdFlags|log.Lshortfile)

	c := cron.New()

	// Bind cron start/stop ke lifecycle FX
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			cronLogger.Println("Starting Transaction Job...")
			fmt.Fprintf(gin.DefaultWriter, "[CRON] Starting Transaction Job...\n")

			// Register cron job
			_, err := c.AddFunc("0 1 * * *", func() {
				cronLogger.Println("Transaction job execution started")

				req := params.CreateTransactionRequest{
					Date: time.Now().Format("2006-01-02"),
				}

				result, err := p.TransactionSvc.Create(context.Background(), req)
				if err != nil {
					cronLogger.Printf("ERROR: Transaction job failed: %v\n", err)
					fmt.Fprintf(gin.DefaultWriter, "[CRON][ERROR] Transaction job failed: %v\n", err)
					return
				}

				cronLogger.Printf("SUCCESS: Transaction job executed. Created %d transaction(s)\n", len(result))
				fmt.Fprintf(gin.DefaultWriter, "[CRON] Transaction job executed successfully\n")
			})
			if err != nil {
				cronLogger.Printf("ERROR: Failed to register cron job: %v\n", err)
				return err
			}

			// Start cron
			c.Start()
			cronLogger.Println("Transaction Job scheduler started")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cronLogger.Println("Stopping Transaction Job...")
			fmt.Fprintf(gin.DefaultWriter, "[CRON] Stopping Transaction Job...\n")
			c.Stop()
			cronLogger.Println("Transaction Job stopped gracefully")

			// Close log file on shutdown
			if err := logFile.Close(); err != nil {
				fmt.Fprintf(gin.DefaultWriter, "[CRON][ERROR] Failed to close log file: %v\n", err)
			}
			return nil
		},
	})
}
