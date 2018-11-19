package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d2r2/go-i2c"
)

var output1 uint8
var output2 uint8
var myExitflag bool = false

func slaveWriteRegU8(slave uint8, addr uint8, data uint8) {

	i2c, err := i2c.NewI2C(slave, 1)

	if err != nil {
		log.Fatal(err)
	}
	// Free I2C connection on exit
	defer i2c.Close()

	err = i2c.WriteRegU8(addr, data)
	if err != nil {
		log.Fatal(err)

	}
	i2c.Close()
}

func card1Test() {

	if output1 < 0x10 {
		output1 = 0x10
	}
	slaveWriteRegU8(0x20, 0, 0)
	slaveWriteRegU8(0x20, 0x9, output1)
	output1 = output1 + 0x10

}

func card2Test() {

	if output2 < 0x10 {
		output2 = 0x10
	}
	slaveWriteRegU8(0x21, 0, 0)
	slaveWriteRegU8(0x21, 0x9, output2)
	output2 = output2 + 0x10

}

func main() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		myExitflag = true
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 4 second to finish processing")

		time.Sleep(4 * time.Second)
		os.Exit(0)

	}()

	for {
		card1Test()
		time.Sleep(time.Millisecond * 2000)
		card2Test()
		time.Sleep(time.Millisecond * 2000)
		if myExitflag {
			break
		}

	}
	slaveWriteRegU8(0x20, 0x9, 0x00) // turn off outputs on card 1
	slaveWriteRegU8(0x21, 0x9, 0x00) // turn off outputs on card 2

}
