# bme280

This is a simple driver for the Bosch Sensortec bme280 environment sensor.
The datasheet can be found here:
https://ae-bst.resource.bosch.com/media/_tech/media/datasheets/BST-BME280_DS001-11.pdf

ATM the driver is expecting the i2c-bus. SPI-Bus could be added in the future.

# Usage
For Normal mode and default values you can use like that:

	package main
	
	import "log"
	import "github.com/davecheney/i2c"
	import "github.com/quinte17/bme280" 
	
	func main() {
		dev, err := i2c.New(0x77, 1)
		if err != nil {
			log.Print(err)
		}
		bme, err := bme280.NewI2CDriver(dev)
		if err != nil {
			log.Print(err)
		}
		log.Print(bme.Readenv())
	}

If you want to use the forced mode use like that:

	package main
	
	import "log"
	import "github.com/davecheney/i2c"
	import "github.com/quinte17/bme280" 
	
	func main() {
		dev, err := i2c.New(0x77, 1)
		if err != nil {
			log.Print(err)
		}
		bme, err := bme280.NewI2CDriver(dev, bme280.OptHumOversampling(1), bme280.OptTempOversampling(1), bme280.OptPressOversampling(1))
		if err != nil {
			log.Print(err)
		}
		bme.Option(bme280.OptMode("forced")) // start one measurement
		bme.WaitForMeasurement()             // sleep long enough so the measurement will be finished.
		log.Print(bme.Readenv())
	}
