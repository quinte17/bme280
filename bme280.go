package bme280

import "github.com/davecheney/i2c"
import "time"

type BME280 struct {
	i2c   *i2c.I2C
	calib struct {
		temp struct {
			T1 uint16
			T2 int16
			T3 int16
		}
		press struct {
			P1 uint16
			P2 int16
			P3 int16
			P4 int16
			P5 int16
			P6 int16
			P7 int16
			P8 int16
			P9 int16
		}
		hum struct {
			H1 uint8
			H2 int16
			H3 uint8
			H4 int16
			H5 int16
			H6 int8
		}
	}
	raw [8]byte
}

func (bme *BME280) read(reg byte, data []byte) (int, error) {
	// first we have to write register adress
	_, err := bme.i2c.Write([]byte{reg})
	if err != nil {
		return 0, err
	}
	// now we can read the data
	return bme.i2c.Read(data)
}

func (bme *BME280) bootFinished() (err error) {
	var x [1]byte
	for x[0] != 0x60 {
		_, err = bme.read(REG_id, x[:])
		time.Sleep(50 * time.Millisecond)
	}
	return err
}

func (bme *BME280) readCalibdata() (err error) {
	// read calibration data
	var calib1 [26]byte
	var calib2 [16]byte
	bme.read(REG_calib00, calib1[:])
	bme.read(REG_calib26, calib2[:])

	convert(calib1[0:6], &bme.calib.temp)
	convert(calib1[6:24], &bme.calib.press)

	convert(calib1[25:], &bme.calib.hum.H1)
	convert(calib2[0:2], &bme.calib.hum.H2)
	convert(calib2[2:3], &bme.calib.hum.H3)
	// H4 and H5 are a little bit tricky alligned.
	bme.calib.hum.H4 = int16(calib2[3])<<4 | int16(calib2[4]&0x0F)
	bme.calib.hum.H5 = int16(calib2[5])<<4 | int16(calib2[4]&0xF0)>>4
	convert(calib2[6:], &bme.calib.hum.H6)

	return err
}

func (bme *BME280) initialize() (err error) {
	// wait for finished initialisation
	bme.bootFinished()
	// get calibrationdata
	bme.readCalibdata()

	// initialize bme
	bme.i2c.Write([]byte{REG_ctrl_hum, OPT_hum_oversampling_x1})
	bme.i2c.Write([]byte{REG_ctrl_meas, OPT_temp_oversampling_x1 | OPT_press_oversampling_x1 | OPT_mode_normal})
	bme.i2c.Write([]byte{REG_config, OPT_config_standbytime_1000})

	return err
}

// latch all data in
func (bme *BME280) ReadRaw() (err error) {
	_, err = bme.read(REG_press_msb, bme.raw[:])
	return err
}

func New(i2c *i2c.I2C) (*BME280, error) {
	bme := BME280{
		i2c: i2c,
	}

	bme.initialize()

	return &bme, nil
}
