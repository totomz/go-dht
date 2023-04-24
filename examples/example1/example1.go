package main

import (
	"github.com/totomz/go-dht"
	"log"
)

func main() {

	println("***************************************************************************************************")
	println("*** You can change verbosity of output, to modify logging level of module \"dht\"")
	println("*** Uncomment/comment corresponding lines with call to ChangePackageLogLevel(...)")
	println("***************************************************************************************************")
	// Uncomment/comment next line to suppress/increase verbosity of output
	// logger.ChangePackageLogLevel("dht", logger.InfoLevel)

	// sensorType := dht.DHT11
	sensorType := dht.AM2302
	// sensorType := dht.DHT12
	// Read DHT11 sensor data from specific pin, retrying 10 times in case of failure.
	pin := 1
	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(sensorType, pin, false, 10)
	if err != nil {
		log.Fatal(err)
	}
	// print temperature and humidity
	log.Printf("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)", sensorType, temperature, humidity, retried)
}
