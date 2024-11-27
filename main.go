package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"olympos.io/encoding/edn"
)

type Config struct {
	CfgVers      int    `edn:"meta/version"`
	PreferredFmt string `edn:"preferred-format"`
	FileNameFmt  string `edn:"journal/file-name-format"`
}

func main() {
	// Define command line flags
	lsqDirName := flag.String("d", "Logseq", "The main Logseq directory to use.")
	lsqCfgDirName := flag.String("l", "logseq", "The Logseq configuration directory to use.")
	lsqCfgFileName := flag.String("c", "config.edn", "The config.edn file to use.")
	editor := flag.String("e", "EDITOR", "The editor to use.")

	// Parse flags
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Construct paths
	lsqDir := filepath.Join(homeDir, *lsqDirName)
	lsqCfgDir := filepath.Join(lsqDir, *lsqCfgDirName)
	cfgFile := filepath.Join(lsqCfgDir, *lsqCfgFileName)

	// Read config file to determine preferred format
	configData, err := os.ReadFile(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		os.Exit(1)
	}

	// Set defaults before extracting data from config file:
	var cfg = &Config{
		CfgVers:      1,
		PreferredFmt: "Markdown",
		FileNameFmt:  "yyyy_MM_dd",
	}

	err = edn.Unmarshal(configData, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling config data: %v\n", err)
		os.Exit(1)
	}

	// Construct journals directory path
	journalsDir := filepath.Join(lsqDir, "journals")

	// Create journals directory if it doesn't exist
	err = os.MkdirAll(journalsDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating journals directory: %v\n", err)
		os.Exit(1)
	}

	// Construct today's journal file path
	extension := ".md"
	if cfg.PreferredFmt == "Org" {
		extension = ".org"
	}

	// Get today's date in YYYY_MM_DD format
	today := time.Now().Format(convertDateFormat(cfg.FileNameFmt))

	journalPath := filepath.Join(journalsDir, today+extension)

	// Create file if it doesn't exist
	_, err = os.Stat(journalPath)

	if os.IsNotExist(err) {
		err := os.WriteFile(journalPath, []byte(""), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating journal file: %v\n", err)
			os.Exit(1)
		}
	}

	// Get editor from environment
	editing := selectEditor(*editor)

	// Open file in editor
	cmd := exec.Command(editing, journalPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func selectEditor(editor string) string {
	if editor == "" {
		return "vim"
	}

	checkEnv := os.Getenv(editor)
	if checkEnv != "" {
		return checkEnv
	}

	// Return whatever if not Env Var
	return editor
}

func convertDateFormat(lsqFormat string) string {
	// Map of lsq date format tokens to Go date format
	formatMap := map[string]string{
		"yyyy": "2006",
		"yy":   "06",
		"MM":   "01",
		"M":    "1",
		"dd":   "02",
		"d":    "2",
	}

	goFormat := lsqFormat
	for lsqToken, goToken := range formatMap {
		goFormat = strings.ReplaceAll(goFormat, lsqToken, goToken)
	}

	return goFormat
}
