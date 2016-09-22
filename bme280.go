package bme280

import "io"
import "time"

type BME280 struct {
	i2c   io.ReadWriter
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

type option func(*BME280) error

// Option sets the options specified.
// It returns an option to restore the last arg's previous value.
func (bme *BME280) Option(opts ...option) (err error) {
	for _, opt := range opts {
		err = opt(bme)
		if err != nil {
			return err
		}
	}
	return nil
}

func OptHumOversampling(sampling int) option {
	return func(bme *BME280) error {
		// read
		var r byte
		_, err := bme.read(REG_ctrl_hum, []byte{r})
		if err != nil {
			return err
		}
		// modify
		r = r & ^byte(OPT_hum_mask)
		switch sampling {
		case 0:
			r = r | OPT_hum_oversampling_skipped
		case 1:
			r = r | OPT_hum_oversampling_x1
		case 2:
			r = r | OPT_hum_oversampling_x2
		case 4:
			r = r | OPT_hum_oversampling_x4
		case 8:
			r = r | OPT_hum_oversampling_x8
		case 16:
			r = r | OPT_hum_oversampling_x16
		}
		// write
		_, err = bme.write(REG_ctrl_hum, []byte{r})
		return err
	}
}

func OptTempOversampling(sampling int) option {
	return func(bme *BME280) error {
		// read
		var r byte
		_, err := bme.read(REG_ctrl_meas, []byte{r})
		if err != nil {
			return err
		}
		// modify
		r = r & ^byte(OPT_temp_mask)
		switch sampling {
		case 0:
			r = r | OPT_temp_oversampling_skipped
		case 1:
			r = r | OPT_temp_oversampling_x1
		case 2:
			r = r | OPT_temp_oversampling_x2
		case 4:
			r = r | OPT_temp_oversampling_x4
		case 8:
			r = r | OPT_temp_oversampling_x8
		case 16:
			r = r | OPT_temp_oversampling_x16
		}
		// write
		_, err = bme.write(REG_ctrl_meas, []byte{r})
		return err
	}
}

func OptPressOversampling(sampling int) option {
	return func(bme *BME280) error {
		// read
		var r byte
		_, err := bme.read(REG_ctrl_meas, []byte{r})
		if err != nil {
			return err
		}
		// modify
		r = r & ^byte(OPT_press_mask)
		switch sampling {
		case 0:
			r = r | OPT_press_oversampling_skipped
		case 1:
			r = r | OPT_press_oversampling_x1
		case 2:
			r = r | OPT_press_oversampling_x2
		case 4:
			r = r | OPT_press_oversampling_x4
		case 8:
			r = r | OPT_press_oversampling_x8
		case 16:
			r = r | OPT_press_oversampling_x16
		}
		// write
		_, err = bme.write(REG_ctrl_meas, []byte{r})
		return err
	}
}

func OptMode(mode string) option {
	return func(bme *BME280) error {
		// read
		var r byte
		_, err := bme.read(REG_ctrl_meas, []byte{r})
		if err != nil {
			return err
		}
		// modify
		r = r & ^byte(OPT_mode_mask)
		switch mode {
		case "sleep":
			r = r | OPT_mode_sleep
		case "forced":
			r = r | OPT_mode_forced
		case "normal":
			r = r | OPT_mode_normal
		}
		// write
		_, err = bme.write(REG_ctrl_meas, []byte{r})
		return err
	}
}

func OptFilter(mode int) option {
	return func(bme *BME280) error {
		// read
		var r byte
		_, err := bme.read(REG_config, []byte{r})
		if err != nil {
			return err
		}
		// modify
		r = r & ^byte(OPT_config_filter_mask)
		switch mode {
		case 0:
			r = r | OPT_config_filter_off
		case 2:
			r = r | OPT_config_filter_2
		case 4:
			r = r | OPT_config_filter_4
		case 8:
			r = r | OPT_config_filter_8
		case 16:
			r = r | OPT_config_filter_16
		}
		// write
		_, err = bme.write(REG_config, []byte{r})
		return err
	}
}

func OptStandbytime(time int) option {
	return func(bme *BME280) error {
		// read
		var r byte
		_, err := bme.read(REG_config, []byte{r})
		if err != nil {
			return err
		}
		// modify
		r = (r & ^byte(OPT_config_standbytime_mask))
		switch time {
		case 5:
			r = r | OPT_config_standbytime_0_5
		case 625:
			r = r | OPT_config_standbytime_62_5
		case 125:
			r = r | OPT_config_standbytime_125
		case 250:
			r = r | OPT_config_standbytime_250
		case 500:
			r = r | OPT_config_standbytime_500
		case 1000:
			r = r | OPT_config_standbytime_1000
		case 10:
			r = r | OPT_config_standbytime_10
		case 20:
			r = r | OPT_config_standbytime_20
		}
		// write
		_, err = bme.write(REG_config, []byte{r})
		return err
	}
}

type Envdata struct {
	Temp  float64 `json:"temp"`
	Press float64 `json:"press"`
	Hum   float64 `json:"hum"`
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

func (bme *BME280) write(reg byte, data []byte) (int, error) {
	var tdata []byte
	tdata = append(tdata, reg)
	tdata = append(tdata, data...)

	return bme.i2c.Write(tdata)
}

func (bme *BME280) bootFinished() (err error) {
	var x [1]byte
	for x[0] != 0x60 && err == nil {
		_, err = bme.read(REG_id, x[:])
		time.Sleep(50 * time.Millisecond)
	}
	return err
}

func (bme *BME280) readCalibdata() (err error) {
	// read calibration data
	var calib1 [26]byte
	var calib2 [16]byte
	_, err = bme.read(REG_calib00, calib1[:])
	if err != nil {
		return err
	}
	_, err = bme.read(REG_calib26, calib2[:])
	if err != nil {
		return err
	}

	type tmpt struct {
		idata []byte
		odata interface{}
	}
	tconvert := []tmpt{
		{calib1[0:6], &bme.calib.temp},
		{calib1[6:24], &bme.calib.press},
		{calib1[25:], &bme.calib.hum.H1},
		{calib2[0:2], &bme.calib.hum.H2},
		{calib2[2:3], &bme.calib.hum.H3},
		{calib2[6:], &bme.calib.hum.H6}}

	for _, value := range tconvert {
		err = convert(value.idata, value.odata)
		if err != nil {
			return err
		}
	}

	// H4 and H5 are a little bit tricky alligned.
	bme.calib.hum.H4 = int16(calib2[3])<<4 | int16(calib2[4]&0x0F)
	bme.calib.hum.H5 = int16(calib2[5])<<4 | int16(calib2[4]&0xF0)>>4

	return err
}

func (bme *BME280) initialize(opts ...option) (err error) {
	// wait for finished initialisation
	err = bme.bootFinished()
	if err != nil {
		return err
	}
	// get calibrationdata
	err = bme.readCalibdata()
	if err != nil {
		return err
	}
	// initialize bme
	if len(opts) == 0 {
		return bme.Option(OptHumOversampling(1), OptTempOversampling(1),
			OptPressOversampling(1), OptMode("normal"), OptStandbytime(1000))
	}
	return bme.Option(opts...)
}

// latch all data in
func (bme *BME280) readRaw() (err error) {
	_, err = bme.read(REG_press_msb, bme.raw[:])
	return err
}

// calculate enviroment data
func (bme *BME280) Readenv() (env Envdata, err error) {
	err = bme.readRaw()
	traw := int32(bme.raw[3])<<12 | int32(bme.raw[4])<<4 | int32(bme.raw[5])>>4
	praw := int32(bme.raw[0])<<12 | int32(bme.raw[1])<<4 | int32(bme.raw[2])>>4
	hraw := int32(bme.raw[6])<<8 | int32(bme.raw[7])

	t, tfine := bme.temp(traw)
	p := bme.press(praw, tfine)
	h := bme.hum(hraw, tfine)

	env.Temp = t
	env.Press = p / 100
	env.Hum = h
	return env, err
}

func (bme *BME280) temp(raw int32) (float64, int32) {
	calt := bme.calib.temp
	var v1, v2, t float64
	var tfine int32
	v1 = (float64(raw)/16384.0 - float64(calt.T1)/1024.0) *
		float64(calt.T2)
	v2 = (float64(raw)/131072.0 - float64(calt.T1)/8192.0) *
		(float64(raw)/131072.0 - float64(calt.T1)/8192.0) *
		float64(calt.T3)
	tfine = int32(v1 + v2)
	t = (v1 + v2) / 5120.0
	return t, tfine
}

func (bme *BME280) press(raw int32, tfine int32) float64 {
	calp := bme.calib.press
	var v1, v2, p float64
	v1 = float64(tfine)/2.0 - 64000.0
	v2 = v1 * v1 * (float64(calp.P6) / 32768.0)
	v2 = v2 + v1*(float64(calp.P5)*2.0)
	v2 = v2/4.0 + (float64(calp.P4) * 65536.0)
	v1 = (float64(calp.P3)*v1*v1/524288.0 + float64(calp.P2)*v1) / 524288.0
	v1 = (1.0 + v1/32768.0) * float64(calp.P1)
	if v1 == 0 {
		return 0
	}
	p = 1048576.0 - float64(raw)
	p = (p - v2/4096.0) * 6250.0 / v1
	v1 = float64(calp.P9) * p * p / 2147483648.0
	v2 = p * float64(calp.P8) / 32768.0
	p = p + (v1+v2+float64(calp.P7))/16.0
	return p
}

func (bme *BME280) hum(raw int32, tfine int32) float64 {
	calh := bme.calib.hum
	var h float64
	h = float64(tfine) - 76800.0
	h = (float64(raw) - float64(calh.H4)*64.0 +
		float64(calh.H5)/16384.0*h) * float64(calh.H2) /
		65536.0 * (1.0 + float64(calh.H6)/67108864.0*h*
		(1.0+float64(calh.H3)/67108864.0*h))
	h = h * (1.0 - float64(calh.H1)*h/524288.0)

	if h > 100.0 {
		h = 100.0
	} else if h < 0.0 {
		h = 0.0
	}
	return h
}

// NewI2CDriver initializes the bme280 device to use the i2c-bus for communication.
// It is expecting the i2c bus as a ReadWriter-Interface.
// If no options given it will be initialized with:
// OptHumOversampling(1), OptTempOversampling(1), OptPressOversampling(1), OptMode("normal"), OptStandbytime(1000)
func NewI2CDriver(i2c io.ReadWriter, ops ...option) (*BME280, error) {
	bme := BME280{
		i2c: i2c,
	}

	return &bme, bme.initialize(ops...)
}
