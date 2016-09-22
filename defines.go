package bme280

const (
	// Register
	reg_id         = 0xD0
	reg_reset      = 0xE0
	reg_ctrl_hum   = 0xF2
	reg_status     = 0xF3
	reg_ctrl_meas  = 0xF4
	reg_config     = 0xF5
	reg_press_msb  = 0xF7
	reg_press_lsb  = 0xF8
	reg_press_xlsb = 0xF9
	reg_temp_msb   = 0xFA
	reg_temp_lsb   = 0xFB
	reg_temp_xlsb  = 0xFC
	reg_hum_msb    = 0xFD
	reg_hum_lsb    = 0xFE

	reg_calib00 = 0x88
	reg_calib26 = 0xE1

	// Options
	opt_press_oversampling_skipped = 0x00
	opt_press_oversampling_x1      = 0x04
	opt_press_oversampling_x2      = 0x08
	opt_press_oversampling_x4      = 0x0C
	opt_press_oversampling_x8      = 0x10
	opt_press_oversampling_x16     = 0x14
	opt_press_mask                 = 0x1C

	opt_temp_oversampling_skipped = 0x00
	opt_temp_oversampling_x1      = 0x20
	opt_temp_oversampling_x2      = 0x40
	opt_temp_oversampling_x4      = 0x60
	opt_temp_oversampling_x8      = 0x80
	opt_temp_oversampling_x16     = 0xA0
	opt_temp_mask                 = 0xE0

	opt_hum_oversampling_skipped = 0x00
	opt_hum_oversampling_x1      = 0x01
	opt_hum_oversampling_x2      = 0x02
	opt_hum_oversampling_x4      = 0x03
	opt_hum_oversampling_x8      = 0x04
	opt_hum_oversampling_x16     = 0x05
	opt_hum_mask                 = 0x07

	opt_mode_sleep  = 0x00
	opt_mode_forced = 0x01
	opt_mode_normal = 0x03
	opt_mode_mask   = 0x03

	opt_config_standbytime_0_5  = 0x00
	opt_config_standbytime_62_5 = 0x20
	opt_config_standbytime_125  = 0x40
	opt_config_standbytime_250  = 0x60
	opt_config_standbytime_500  = 0x80
	opt_config_standbytime_1000 = 0xA0
	opt_config_standbytime_10   = 0xC0
	opt_config_standbytime_20   = 0xE0
	opt_config_standbytime_mask = 0xE0

	opt_config_filter_off  = 0x00
	opt_config_filter_2    = 0x04
	opt_config_filter_4    = 0x08
	opt_config_filter_8    = 0x0C
	opt_config_filter_16   = 0x10
	opt_config_filter_mask = 0x1C

	opt_config_enable_3wire = 0x01
)
