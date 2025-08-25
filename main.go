package main

import (
	"ddnsd/config"
	"ddnsd/internal"
	"ddnsd/utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		utils.LogWarning("Failed to load .env file: %v", err)
	}

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.LogError("Configuration error: %v", err)
		os.Exit(1)
	}

	// Initialize DNS provider
	provider, err := internal.NewDNSProvider(cfg)
	if err != nil {
		utils.LogError("Failed to initialize DNS provider: %v", err)
		os.Exit(1)
	}

	// Print configuration summary
	config.PrintConfigSummary(cfg)

	// Run initial update
	utils.LogInfo("Starting initial update...")
	internal.RunSequentialUpdates(provider, cfg)

	// Set up scheduled updates
	scheduler := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	spec := fmt.Sprintf("@every %ds", cfg.Interval)

	_, err = scheduler.AddFunc(spec, func() {
		internal.RunSequentialUpdates(provider, cfg)
	})
	if err != nil {
		utils.LogError("Failed to set up scheduler: %v", err)
		os.Exit(1)
	}

	scheduler.Start()
	defer scheduler.Stop()

	utils.LogInfo("DDNS service started successfully. Press Ctrl+C to exit.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	utils.LogInfo("Shutting down DDNS service...")
}