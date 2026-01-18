.DEFAULT_GOAL := cli

# CLI-only build (no OpenGL/CGO dependencies)
cli:
	go build -o ffvi_editor_cli.exe .\main_cli.go

# Test CLI build
test-cli:
	go run .\main_cli.go combat-pack --mode smoke

# Full GUI build (requires MinGW/CGO)
gui:
	set CGO_ENABLED=1 && go build -ldflags "-s -H=windowsgui" -o "FFVIPR_Save_Editor.exe"

# Original build target (requires UPX and CGO)
build:
	go build -ldflags "-s -H=windowsgui" -o "FFVIPR_Save_Editor.exe"
	upx -9 -k "FFVIPR_Save_Editor.exe"
	rm "FFVIPR_Save_Editor.ex~"

# E2E test for Combat Pack
test-combat-pack:
	go run test_combat_pack_standalone.go

# Setup tools
setup:
	go install github.com/tc-hib/go-winres@latest
	go-winres simply --icon icon.png

# Clean build artifacts
clean:
	del /F /Q ffvi_editor_cli.exe FFVIPR_Save_Editor.exe 2>nul || echo Clean complete

# Help
help:
	@echo FF6 Save Editor - Build Targets
	@echo.
	@echo   make cli              - Build CLI-only version (no GUI, no CGO)
	@echo   make test-cli         - Test CLI build with smoke tests
	@echo   make gui              - Build full GUI version (requires MinGW)
	@echo   make build            - Original build with compression
	@echo   make test-combat-pack - Run E2E tests
	@echo   make setup            - Install build tools
	@echo   make clean            - Remove build artifacts
	@echo   make help             - Show this help
