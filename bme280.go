package bme280

import "github.com/davecheney/i2c"

//import "fmt"
import "encoding/binary"
import "bytes"
import "time"

const (
	// Register
	REG_id         = 0xD0
	REG_reset      = 0xE0
	REG_ctrl_hum   = 0xF2
	REG_status     = 0xF3
	REG_ctrl_meas  = 0xF4
	REG_config     = 0xF5
	REG_press_msb  = 0xF7
	REG_press_lsb  = 0xF8
	REG_press_xlsb = 0xF9
	REG_temp_msb   = 0xFA
	REG_temp_lsb   = 0xFB
	REG_temp_xlsb  = 0xFC
	REG_hum_msb    = 0xFD
	REG_hum_lsb    = 0xFE

	REG_calib00 = 0x88
	REG_calib26 = 0xE1

	// Options
	OPT_press_oversampling_skipped = 0x00
	OPT_press_oversampling_x1      = 0x04
	OPT_press_oversampling_x2      = 0x08
	OPT_press_oversampling_x4      = 0x0C
	OPT_press_oversampling_x8      = 0x10
	OPT_press_oversampling_x16     = 0x14

	OPT_temp_oversampling_skipped = 0x00
	OPT_temp_oversampling_x1      = 0x20
	OPT_temp_oversampling_x2      = 0x40
	OPT_temp_oversampling_x4      = 0x60
	OPT_temp_oversampling_x8      = 0x80
	OPT_temp_oversampling_x16     = 0xA0

	OPT_hum_oversampling_skipped = 0x00
	OPT_hum_oversampling_x1      = 0x01
	OPT_hum_oversampling_x2      = 0x02
	OPT_hum_oversampling_x4      = 0x03
	OPT_hum_oversampling_x8      = 0x04
	OPT_hum_oversampling_x16     = 0x05

	OPT_mode_sleep  = 0x00
	OPT_mode_forced = 0x01
	OPT_mode_normal = 0x03

	OPT_config_standbytime_0_5  = 0x00
	OPT_config_standbytime_62_5 = 0x20
	OPT_config_standbytime_125  = 0x40
	OPT_config_standbytime_250  = 0x60
	OPT_config_standbytime_500  = 0x80
	OPT_config_standbytime_1000 = 0xA0
	OPT_config_standbytime_10   = 0xC0
	OPT_config_standbytime_20   = 0xE0

	OPT_config_filter_off = 0x00
	OPT_config_filter_2   = 0x04
	OPT_config_filter_4   = 0x08
	OPT_config_filter_8   = 0x0C
	OPT_config_filter_16  = 0x10

	OPT_config_enable_3wire = 0x01
)

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

func convert(b []byte, data interface{}) error {
	buf := bytes.NewReader(b)
	return binary.Read(buf, binary.LittleEndian, data)
}

func NewBME280(i2c *i2c.I2C) (*BME280, error) {
	bme := BME280{
		i2c: i2c,
	}
	// initialize bme
	//	bme.i2c.Write([]byte{REG_ctrl_hum, OPT_hum_oversampling_x1})
	//	bme.i2c.Write([]byte{REG_ctrl_meas, OPT_temp_oversampling_x1 | OPT_press_oversampling_x1 | OPT_mode_normal})
	//	bme.i2c.Write([]byte{REG_config, OPT_config_standbytime_1000})

	// read some data
	var x [1]byte
	for x[0] != 0x60 {
		bme.read(REG_id, x[:])
		time.Sleep(50 * time.Millisecond)
	}

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
	//convert(calib2[3:5], &bme.calib.hum.H4)
	//convert(calib2[4:6], &bme.calib.hum.H5)
	convert(calib2[6:], &bme.calib.hum.H6)

	return &bme, nil
}
