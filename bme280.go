package bme280

import "github.com/davecheney/i2c"

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
	// Bitnames

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
	i2c *i2c.I2C
}

func NewBME280(i2c *i2c.I2C) (*BME280, error) {
	bme := BME280{
		i2c: i2c,
	}
	// initialize bme
	bme.i2c.Write([]byte{REG_ctrl_hum, OPT_hum_oversampling_x1})
	bme.i2c.Write([]byte{REG_ctrl_meas, OPT_temp_oversampling_x1 | OPT_press_oversampling_x1 | OPT_mode_normal})
	bme.i2c.Write([]byte{REG_config, OPT_config_standbytime_1000})

	return &bme, nil
}
