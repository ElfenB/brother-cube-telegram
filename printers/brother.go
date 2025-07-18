package printers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Printer struct{}

// Creates a new Printer instance and prints its version and info
// Can be used to initialize the printer and check its status
func NewPrinter() *Printer {
	printer := &Printer{}

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

func (p *Printer) PrintLabelYolo(label string) error {
	cmd := exec.Command("ptouch-print", "--text", label, "--fontsize", "64")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error printing label: %v, output: %s", err, output)
	}

	log.Printf("Label printed successfully: %s\n", label)
	return nil
}

func (p *Printer) PreviewLabel(label string) ([]byte, error) {
	cmd := exec.Command("ptouch-print", "--text", label, "--writepng", "draft.png")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error previewing label: %v, output: %s", err, output)
	}

	path := os.Getenv("HOME") + "/draft.png"
	fileContent, _ := os.ReadFile(path)

	log.Printf("Label previewed successfully: %s\n", label)
	return fileContent, nil
}
