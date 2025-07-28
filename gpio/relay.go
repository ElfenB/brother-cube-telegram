package gpio

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	turnOffAction = "Turn OFF"
	turnOnAction  = "Turn ON"
)

// GPIO relay controller
type Relay struct {
	pin       gpio.PinIO
	pinNum    int
	lastState gpio.Level
}

// Creates a new relay controller for the specified GPIO pin
func NewRelay(pinNumber int) (*Relay, error) {
	// Initialize periph.io host
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize host: %v", err)
	}

	// Get the GPIO pin
	pin := gpioreg.ByName(fmt.Sprintf("GPIO%d", pinNumber))
	if pin == nil {
		return nil, fmt.Errorf("failed to find GPIO%d", pinNumber)
	}

	log.Printf("GPIO%d initialized as relay control", pinNumber)

	return &Relay{
		pin:       pin,
		pinNum:    pinNumber,
		lastState: gpio.High,
	}, nil
}

// Activates the relay (sets GPIO Low)
func (r *Relay) TurnOn() error {
	if err := r.pin.Out(gpio.Low); err != nil {
		return fmt.Errorf("failed to turn on relay on GPIO%d: %v", r.pinNum, err)
	}
	r.lastState = gpio.Low
	log.Printf("Relay on GPIO%d turned ON", r.pinNum)
	return nil
}

// Deactivates the relay (sets GPIO High)
func (r *Relay) TurnOff() error {
	if err := r.pin.Out(gpio.High); err != nil {
		return fmt.Errorf("failed to turn off relay on GPIO%d: %v", r.pinNum, err)
	}
	r.lastState = gpio.High
	log.Printf("Relay on GPIO%d turned OFF", r.pinNum)
	return nil
}

// Toggles the relay state
func (r *Relay) Toggle() error {
	if r.lastState == gpio.Low {
		return r.TurnOff()
	}
	return r.TurnOn()
}

// Returns the current state of the relay
func (r *Relay) GetState() bool {
	return r.lastState == gpio.Low
}

// Cleans up the GPIO pin
func (r *Relay) Close() error {
	// Turn off relay before closing
	if err := r.TurnOff(); err != nil {
		log.Printf("Warning: failed to turn off relay during close: %v", err)
	}

	log.Printf("GPIO%d relay closed", r.pinNum)
	return nil
}
