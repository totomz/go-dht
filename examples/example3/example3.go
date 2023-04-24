package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/d2r2/go-shell"
	"github.com/totomz/go-dht"
)

func main() {

	log.Println("***************************************************************************************************")
	log.Println("*** You can change verbosity of output, to modify logging level of module \"dht\"")
	log.Println("*** Uncomment/comment corresponding lines with call to ChangePackageLogLevel(...)")
	log.Println("***************************************************************************************************")
	log.Println("*** Massive stress test of sensor reading, printing in the end summary statistical results")
	log.Println("***************************************************************************************************")

	// create context with cancellation possibility
	ctx, cancel := context.WithCancel(context.Background())
	// use done channel as a trigger to exit from signal waiting goroutine
	done := make(chan struct{})
	defer close(done)
	// build actual signal list to control
	signals := []os.Signal{os.Kill, os.Interrupt}
	if shell.IsLinuxMacOSFreeBSD() {
		signals = append(signals, syscall.SIGTERM)
	}
	// run goroutine waiting for OS termination events, including keyboard Ctrl+C
	shell.CloseContextOnSignals(cancel, done, signals...)

	// sensorType := dht.DHT11
	// sensorType := dht.AM2302
	sensorType := dht.DHT12
	pin := 1
	totalRetried := 0
	totalMeasured := 0
	totalFailed := 0
	term := false
	for i := 0; i < 300; i++ {
		// Read DHT11 sensor data from specific pin, retrying 10 times in case of failure.
		temperature, humidity, retried, err :=
			dht.ReadDHTxxWithContextAndRetry(ctx, sensorType, pin, false, 10)
		totalMeasured++
		totalRetried += retried
		if err != nil && ctx.Err() == nil {
			totalFailed++
			log.Printf("error: %v", err)
			continue
		}
		// print temperature and humidity
		if ctx.Err() == nil {
			log.Printf("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)", sensorType, temperature, humidity, retried)
		}
		select {
		// Check for termination request.
		case <-ctx.Done():
			log.Printf("error: Termination pending: %s", ctx.Err())
			term = true
			// sleep 1.5-2 sec before next round
			// (recommended by specification as "collecting period")
		case <-time.After(2000 * time.Millisecond):
		}
		if term {
			break
		}
	}
	log.Println("====================================================================")
	log.Printf("Total measured = %v, total retried = %v, total failed = %v", totalMeasured, totalRetried, totalFailed)
	log.Println("====================================================================")
}
