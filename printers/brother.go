package printers

import (
	"fmt"
	"os/exec"
)

type Printer struct{}

// Creates a new Printer instance and prints its version and info
// Can be used to initialize the printer and check its status
func NewPrinter() *Printer {
	printer := &Printer{}

	version := printer.GetVersion()
	fmt.Printf("Printer version: %s\n", version)

	info, err := printer.GetPrinterInfo()
	if err != nil {
		fmt.Printf("Error getting printer info: %v\n", err)
		return nil
	}

	fmt.Printf("Printer info: %s\n", info)
	return printer
}

// Returns the printer driver version
func (p *Printer) GetVersion() string {
	cmd := exec.Command("ptouch-print", "--version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "Unknown version"
	}

	return string(output)
}

// Returns information about the printer
func (p *Printer) GetPrinterInfo() (string, error) {
	cmd := exec.Command("ptouch-print", "--info")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
