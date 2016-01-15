# bme280

This is a simple driver for the Bosch Sensortec bme280 environment sensor.
The datasheet can be found here:
http://ae-bst.resource.bosch.com/media/products/dokumente/bme280/BST-BME280_DS001-11.pdf

ATM the driver is expecting the i2c-bus.

# Usage
	import "log"
	import "github.com/davecheney/i2c"
	import "github.com/quinte17/bme280" 
	
	func main() {
		dev, err := i2c.New(0x77, 1)
		if err != nil {
			log.Print(err)
		}
		bme, err := bme280.New(dev)
		if err != nil {
			log.Print(err)
		}
		log.Print(bme.Readenv())
	}
