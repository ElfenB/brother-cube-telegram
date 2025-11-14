package printers

import (
	"brother-cube-telegram/config"
	"brother-cube-telegram/gpio"
	"brother-cube-telegram/logger"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

const print = "ptouch-print"

type Printer struct {
	config        *config.Config
	relay         *gpio.Relay
	shutdownTimer *time.Timer
	timerMutex    sync.Mutex
}

// Creates a new Printer instance and prints its version and info
// Can be used to initialize the printer and check its status
// Takes an optional relay for automatic printer power management
func NewPrinter(relay *gpio.Relay) *Printer {
	cfg := config.Get()

	printer := &Printer{
		config:        cfg,
		relay:         relay,
		shutdownTimer: time.NewTimer(cfg.Printer.GetAutoShutdownDelay()),
	}

	// Stop the initial timer since we haven't started the printer yet
	printer.shutdownTimer.Stop()

	// Start the auto-shutdown goroutine if relay is available
	if relay != nil {
		go printer.autoShutdownRoutine()
	}

	// First ensure printer is on, then get version and info
	if err := printer.ensurePrinterOn(); err != nil {
		logger.Warn("Could not ensure printer is on during initialization: %v", err)
	}

	version := printer.GetVersion()
	logger.Info("Printer version: %s", version)

	info, err := printer.GetPrinterInfo()
	if err != nil {
		logger.Error("Error getting printer info: %v", err)
		return nil
	}

	logger.Info("Printer info: %s", info)
	return printer
}

// Runs in a goroutine and handles automatic printer shutdown
func (p *Printer) autoShutdownRoutine() {
	for {
		<-p.shutdownTimer.C
		p.timerMutex.Lock()
		if p.relay != nil && p.relay.GetState() {
			logger.Info("Auto-shutdown: Turning off printer after %d minutes of inactivity", p.config.Printer.AutoShutdownDelayMinutes)
			if err := p.relay.TurnOff(); err != nil {
				logger.Error("Error during auto-shutdown: %v", err)
			}
		}
		p.timerMutex.Unlock()
	}
}

// Resets the auto-shutdown countdown timer
func (p *Printer) resetAutoShutdownTimer() {
	if p.relay == nil {
		return // No auto-shutdown if no relay
	}

	p.timerMutex.Lock()
	defer p.timerMutex.Unlock()

	// Stop the current timer and reset it
	if !p.shutdownTimer.Stop() {
		// If timer already fired, drain the channel
		select {
		case <-p.shutdownTimer.C:
		default:
		}
	}
	p.shutdownTimer.Reset(p.config.Printer.GetAutoShutdownDelay())
}

// Ensures the printer is powered on via the relay if available
func (p *Printer) ensurePrinterOn() error {
	if p.relay == nil {
		// No relay available, assume printer is always on
		return nil
	}

	if p.relay.GetState() {
		// Relay is already on, printer should be powered
		// Reset the auto-shutdown timer since we're using the printer
		p.resetAutoShutdownTimer()
	} else {
		err := p.relay.TurnOn()
		if err != nil {
			return fmt.Errorf("failed to turn on printer via relay: %v", err)
		}
	}

	// Check via the printer info command if it responds
	_, err := p.execDirect(infoCmdArg)
	if err != nil {
		// Retry with increasing delay
		for i := range p.config.Printer.RetryAttempts - 1 {
			time.Sleep(p.config.Printer.GetRetryDelay(i))
			_, err = p.execDirect(infoCmdArg)
			if err == nil {
				break // Successfully powered on
			}
		}
		if err != nil {
			// If still not responding, return an error
			return fmt.Errorf("printer did not respond after turning on: %v", err)
		}
	}

	// Start the auto-shutdown timer since we just turned on the printer
	p.resetAutoShutdownTimer()
	return nil
}

// Returns the printer driver version
func (p *Printer) GetVersion() string {
	output, err := p.exec("--version")

	if err != nil {
		return "Unknown version"
	}

	return output
}

// Returns information about the printer
func (p *Printer) GetPrinterInfo() (string, error) {
	output, err := p.exec(infoCmdArg)

	if err != nil {
		return "", err
	}

	return output, nil
}

func (p *Printer) PrintLabelYolo(label string) error {
	fontSize := fmt.Sprintf("%d", p.config.Printer.FontSize)
	output, err := p.exec(fontSizeCmdArg, fontSize, textCmdArg, label)

	if err != nil {
		return fmt.Errorf("error printing label: %v, output: %s", err, output)
	}

	logger.Info("Label printed successfully: %s", label)
	return nil
}

func (p *Printer) PrintLabel(label string, fontSize int) error {
	fontSizeStr := fmt.Sprintf("%d", fontSize)
	output, err := p.exec(fontSizeCmdArg, fontSizeStr, textCmdArg, label)

	if err != nil {
		return fmt.Errorf("error printing label: %v, output: %s", err, output)
	}

	logger.Info("Label printed successfully: %s", label)
	return nil
}

func (p *Printer) PrintLabelWithPreset(label string, preset *config.Preset) error {
	fontSizeStr := fmt.Sprintf("%d", preset.FontSize)

	var args []string
	if preset.FontFamily != "" {
		args = append(args, fontCmdArg, preset.FontFamily)
	}
	args = append(args, fontSizeCmdArg, fontSizeStr, textCmdArg, label)

	output, err := p.exec(args...)

	if err != nil {
		return fmt.Errorf("error printing label with font: %v, output: %s", err, output)
	}

	logger.Info("Label printed successfully with font '%s': %s", preset.FontFamily, label)
	return nil
}

func (p *Printer) PreviewLabel(label string, userIdent int64) ([]byte, error) {
	draftsFolder := p.config.Printer.DraftsFolder

	// Ensure the drafts folder exists
	if _, err := os.Stat(draftsFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(draftsFolder, p.config.Printer.GetFolderPermissions()); err != nil {
			return nil, fmt.Errorf("failed to create drafts folder: %v", err)
		}
		logger.Info("Created drafts folder: %s", draftsFolder)
	}

	// Construct the filePath name based on user identifier (e.g. draft-23479234.png)
	filePath := fmt.Sprintf("%s/draft-%d.png", draftsFolder, userIdent)

	output, err := p.exec(fontSizeCmdArg, fmt.Sprintf("%d", p.config.Printer.FontSize), textCmdArg, label, writePngCmdArg, filePath)

	if err != nil {
		return nil, fmt.Errorf("error previewing label: %v, output: %s", err, output)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading label preview file: %v", err)
	}

	logger.Info("Label previewed successfully: %s", label)
	return fileContent, nil
}

// Shuts down the printer, stopping the auto-shutdown timer, optionally turning off the relay
func (p *Printer) Close() error {
	if p.relay == nil {
		return nil
	}

	p.timerMutex.Lock()
	defer p.timerMutex.Unlock()

	// Stop the auto-shutdown timer
	if p.shutdownTimer != nil {
		p.shutdownTimer.Stop()
	}

	// Turn off the printer if it's on
	if p.relay.GetState() {
		if err := p.relay.TurnOff(); err != nil {
			logger.Error("Error turning off printer during close: %v", err)
			return err
		}
		logger.Info("Printer turned off during shutdown")
	}

	return nil
}

// Executes a command on the printer and returns the output
// Returns an error if the printer is not powered on
func (p *Printer) exec(arg ...string) (string, error) {
	if err := p.ensurePrinterOn(); err != nil {
		return "", fmt.Errorf("failed to ensure printer is on: %v", err)
	}

	// Log the command being executed
	logger.Debug("Executing command: %s %v", print, arg)

	command := exec.Command(print, arg...)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing command '%s': %v, output: %s", arg, err, output)
	}

	return string(output), nil
}

// Executes a command on the printer directly without ensuring power-on
// Used internally to avoid recursion when checking printer status
func (p *Printer) execDirect(arg ...string) (string, error) {
	// Log the command being executed
	logger.Debug("Executing command: %s %v", print, arg)

	command := exec.Command(print, arg...)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing command '%v': %v, output: %s", arg, err, output)
	}

	return string(output), nil
}
