package main

import (
	"fmt"
	"os"

	"ffvi_editor/global"
	"ffvi_editor/ui"
	"ffvi_editor/ui/forms/editors"
)

func main() {
	// Setup log file
	if err := global.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log: %v\n", err)
	}
	defer global.CloseLogger()

	defer func() {
		if r := recover(); r != nil {
			global.Log("[PANIC] Fatal error: %v", r)
			crashFile := global.WriteCrashReport(r)
			fmt.Fprintf(os.Stderr, "CRASH: %v\n", r)
			if crashFile != "" {
				fmt.Fprintf(os.Stderr, "Crash report: %s\n", crashFile)
			}
			fmt.Fprintf(os.Stderr, "Log file: %s\n", global.LogPath())
			os.Exit(1)
		}
	}()
	
	global.Log("[MAIN] Starting FFVI Save Editor")
	gui := ui.New()
	editors.CreateTextBoxes()
	gui.Load()
	gui.Run()
	gui.App().Quit()
	global.Log("[MAIN] Exiting normally")
}
