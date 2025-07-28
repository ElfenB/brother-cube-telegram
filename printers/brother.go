package printers

import (
	"brother-cube-telegram/gpio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

const print = "ptouch-print"
const printerPowerErrorMsg = "failed to ensure printer is on: %v"
const autoShutdownDelay = 2 * time.Minute
const retryConnectionAttempts = 5

type Printer struct {
	relay         *gpio.Relay
	shutdownTimer *time.Timer
	timerMutex    sync.Mutex
}

// Creates a new Printer instance and prints its version and info
// Can be used to initialize the printer and check its status
// Takes an optional relay for automatic printer power management
func NewPrinter(relay *gpio.Relay) *Printer {
	printer := &Printer{
		relay:         relay,
		shutdownTimer: time.NewTimer(autoShutdownDelay),
	}

	// Stop the initial timer since we haven't started the printer yet
	printer.shutdownTimer.Stop()

	// Start the auto-shutdown goroutine if relay is available
	if relay != nil {
		go printer.autoShutdownRoutine()
	}

	// First ensure printer is on, then get version and info
	if err := printer.ensurePrinterOn(); err != nil {
		log.Printf("Warning: Could not ensure printer is on during initialization: %v", err)
	}

	version := printer.GetVersion()
	log.Printf("Printer version: %s\n", version)

	info, err := printer.GetPrinterInfo()
	if err != nil {
		log.Printf("Error getting printer info: %v\n", err)
		return nil
	}

	log.Printf("Printer info: %s\n", info)
	return printer
}

// Runs in a goroutine and handles automatic printer shutdown
func (p *Printer) autoShutdownRoutine() {
	for {
		<-p.shutdownTimer.C
		p.timerMutex.Lock()
		if p.relay != nil && p.relay.GetState() {
			log.Println("Auto-shutdown: Turning off printer after 5 minutes of inactivity")
			if err := p.relay.TurnOff(); err != nil {
				log.Printf("Error during auto-shutdown: %v", err)
			}
		}
		p.timerMutex.Unlock()
	}
}

// Resets the 5-minute countdown timer
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
	p.shutdownTimer.Reset(autoShutdownDelay)
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
		for i := range retryConnectionAttempts - 1 {
			time.Sleep(time.Duration(i+5) * time.Second)
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
	output, err := p.exec(textCmdArg, label, fontSizeCmdArg, "64")

	if err != nil {
		return fmt.Errorf("error printing label: %v, output: %s", err, output)
	}

	log.Printf("Label printed successfully: %s\n", label)
	return nil
}

func (p *Printer) PreviewLabel(label string) ([]byte, error) {
	output, err := p.exec(textCmdArg, label, writePngCmdArg, "draft.png")

	if err != nil {
		return nil, fmt.Errorf("error previewing label: %v, output: %s", err, output)
	}

	path := os.Getenv("HOME") + "/draft.png"
	fileContent, _ := os.ReadFile(path)

	log.Printf("Label previewed successfully: %s\n", label)
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
			log.Printf("Error turning off printer during close: %v", err)
			return err
		}
		log.Println("Printer turned off during shutdown")
	}

	return nil
}

// Executes a command on the printer and returns the output
// Returns an error if the printer is not powered on
func (p *Printer) exec(arg ...string) (string, error) {
	if err := p.ensurePrinterOn(); err != nil {
		return "", fmt.Errorf(printerPowerErrorMsg, err)
	}

	// Log the command being executed
	log.Printf("Executing command: %s %v\n", print, arg)

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
	log.Printf("Executing command: %s %v\n", print, arg)

	command := exec.Command(print, arg...)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing command '%v': %v, output: %s", arg, err, output)
	}

	return string(output), nil
}
