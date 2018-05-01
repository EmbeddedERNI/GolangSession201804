package main

import (
	"fmt"
	"log"
	"net/http"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
	"time"
	// layout https://webofthings.org/wp-content/uploads/2016/10/pi-gpio.png
	"periph.io/x/periph/host/rpi"
)

var validPins = map[string]gpio.PinIO{
	"11": rpi.P1_11,
}

type Controller struct {
	handler http.HandlerFunc
}

func InitController() {
	_, err := host.Init()
	if err != nil {
		log.Printf("err = %+v\n", err)
	}
	if !rpi.Present() {
		log.Println("Not runing in raspberry pi")
	}

	rpi.P1_11.Out(gpio.High)
	go func() {
		time.Sleep(500 * time.Millisecond)
		rpi.P1_11.Out(gpio.Low)
	}()

}
func CleanController() {
	for _, v := range validPins {
		v.Out(gpio.Low)
	}
}

func (c Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler.ServeHTTP(w, r)
	log.Println("Check Control Parameters")
	if gpiopin, exist := r.URL.Query()["gpio"]; exist {
		if pin, exist := validPins[gpiopin[0]]; exist {
			if state, exist := r.URL.Query()["state"]; exist {
				switch state[0] {
				case "on":
					pin.Out(gpio.High)
					fmt.Fprintf(w, "Your Are controlling the gpio %s=%s", gpiopin, state)
				case "off":
					pin.Out(gpio.Low)
					fmt.Fprintf(w, "Your Are controlling the gpio %s=%s", gpiopin, state)
				default:
					fmt.Fprintf(w, "State must be on or off")
				}
			}
		} else {
			fmt.Fprintf(w, "Invalid Pin %s", gpiopin)
		}
	}
}
