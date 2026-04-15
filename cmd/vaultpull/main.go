package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/sync"
)

const version = "0.1.0"

func main() {
	var (
		showVersion = flag.Bool("version", false, "print version and exit")
		outputFile  = flag.String("output", ".env", "path to output .env file")
		namespace   = flag.String("namespace", "", "Vault namespace (optional)")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("vaultpull v%s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// CLI flags override environment-derived config values.
	if *outputFile != ".env" {
		cfg.OutputFile = *outputFile
	}
	if *namespace != "" {
		cfg.Namespace = *namespace
	}

	syncer, err := sync.New(cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := syncer.Run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
