package touchscreen

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
)

const FT62XX_ADDR = 0x38
const FT62XX_REG_CTPMVENDID = 0xA8
const FT62XX_REG_CHIP_VENDID = 0xA3
const FT62XX_TOUCH_COUNT_REG = 0x02
const FT62XX_TOUCH_VAL_REG = 0x03
const FT62XX_THRESHOLD_ADDR = 0x80
const FT62XX_GESTURE_ADDR = 0x01

const FT5X06_ADDR = 0x38
const FT5X06_REG_CTPMVENDID = 0xA8
const FT5X06_REG_CHIP_VENDID = 0xA3
const FT5X06_TOUCH_COUNT_REG = 0x02
const FT5X06_TOUCH_VAL_REG = 0x03
const FT5X06_THRESHOLD_ADDR = 0x80
const FT5X06_GESTURE_ADDR = 0x01

type TouchScreen struct {
	Device           *i2c.Device
	Touched          bool
	LastScreenChange time.Time
	Debug            bool
	TType            Touchscreen_Type
}

type TouchPoint struct {
	X int
	Y int
}

type Gesture int

const (
	MOVE_UP         = 0x10
	MOVE_LEFT       = 0x14
	MOVE_DOWN       = 0x18
	MOVE_RIGHT      = 0x1c
	ZOOM_IN         = 0x48
	ZOOM_OUT        = 0x49
	NO_GESTURE      = 0x00
	GESTURE_UNKNOWN = 0xff
)

type Touchscreen_Type int

const (
	FT62XX = 1 + iota
	FT5X06
)

func (t *TouchScreen) Init(tt Touchscreen_Type) error {
	if t.Device != nil {
		return errors.New("Device already initialiased")
	}

	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, FT62XX_ADDR)
	if err != nil {
		return err
	}
	t.Device = d

	t.SetThreshold()

	// Set the defaults
	t.TType = tt
	t.Touched = false
	t.Debug = false

	return nil
}

func (t *TouchScreen) SetThreshold() error {
	threshold_buf := make([]byte, 1)
	threshold_buf[0] = 0x80 // Set default threshold of 50%

	var err error
	if t.TType == FT62XX {
		err = t.Device.WriteReg(FT62XX_THRESHOLD_ADDR, threshold_buf)
	} else if t.TType == FT5X06 {
		err = t.Device.WriteReg(FT5X06_THRESHOLD_ADDR, threshold_buf)
	} else {
		return errors.New("Touchscreen Type is not set or is unrecognized")
	}

	if err != nil {
		return errors.New("Failed to set the threshold")
	}

	return nil
}

func (t *TouchScreen) GetVendorID() (byte, error) {
	buf := make([]byte, 1)
	var err error
	if t.TType == FT62XX {
		err = t.Device.ReadReg(FT62XX_REG_CHIP_VENDID, buf)
	} else if t.TType == FT5X06 {
		err = t.Device.ReadReg(FT5X06_REG_CHIP_VENDID, buf)
	} else {
		return 0, errors.New("Touchscreen Type is not set or is unrecognized")
	}

	if err != nil {
		return buf[0], errors.New("Failed to read in the Vendor Reg")
	}
	return buf[0], nil
}

func (t *TouchScreen) GetCTPMVendorId() (byte, error) {
	buf := make([]byte, 1)
	var err error
	if t.TType == FT62XX {
		err = t.Device.ReadReg(FT62XX_REG_CTPMVENDID, buf)
	} else if t.TType == FT5X06 {
		err = t.Device.ReadReg(FT5X06_REG_CTPMVENDID, buf)
	} else {
		return 0, errors.New("Touchscreen Type is not set or is unrecognized")
	}

	if err != nil {
		return buf[0], errors.New("Failed to read in the Vendor Reg")
	}
	return buf[0], nil
}

// Get the numer of touches the device currently has
// or 0 if we've already sent this one
func (t *TouchScreen) GetTouchesCount() (int, error) {
	if t.Device == nil {
		return 0, errors.New("Touchscreen not initialised yet")
	}

	buf := make([]byte, 1)
	var err error
	if t.TType == FT62XX {
		err = t.Device.ReadReg(FT62XX_TOUCH_COUNT_REG, buf)
	} else if t.TType == FT5X06 {
		err = t.Device.ReadReg(FT5X06_TOUCH_COUNT_REG, buf)
	} else {
		return 0, errors.New("Touchscreen Type is not set or is unrecognized")
	}

	if err != nil {
		return 0, errors.New("Error reading register")
	}

	touchCount := int(buf[0])
	if touchCount > 2 {
		touchCount = 0
	}
	if t.Touched == true && touchCount == 0 {
		t.Touched = false
		return 0, err
	}

	if t.Touched == true {
		// We have already sent this touch - not again please
		return 0, nil
	}

	if touchCount > 0 && touchCount < 3 {
		t.Touched = true
	}
	return touchCount, err
}

func (t *TouchScreen) GetGestureType() (Gesture, error) {
	buf := make([]byte, 1)
	var err error
	if t.TType == FT62XX {
		err = t.Device.ReadReg(FT62XX_GESTURE_ADDR, buf)
	} else if t.TType == FT5X06 {
		err = t.Device.ReadReg(FT5X06_GESTURE_ADDR, buf)
	} else {
		return 0, errors.New("Touchscreen Type is not set or is unrecognized")
	}

	if err != nil {
		return GESTURE_UNKNOWN, err
	}

	if buf[0] == MOVE_UP {
		return MOVE_UP, nil
	} else if buf[0] == MOVE_LEFT {
		return MOVE_LEFT, nil
	} else if buf[0] == MOVE_DOWN {
		return MOVE_DOWN, nil
	} else if buf[0] == MOVE_RIGHT {
		return MOVE_RIGHT, nil
	} else if buf[0] == ZOOM_IN {
		return ZOOM_IN, nil
	} else if buf[0] == ZOOM_OUT {
		return ZOOM_OUT, nil
	} else {
		return GESTURE_UNKNOWN, errors.New("Unknown response from device")
	}
}

// Get the touch point of the first touch
// NOTE: Only supports one touch at the moment
func (t *TouchScreen) GetTouches() (TouchPoint, error) {
	ret := TouchPoint{X: 0, Y: 0}
	buf := make([]byte, 4)

	var err error
	if t.TType == FT62XX {
		err = t.Device.ReadReg(FT62XX_TOUCH_VAL_REG, buf)
	} else if t.TType == FT5X06 {
		err = t.Device.ReadReg(FT5X06_TOUCH_VAL_REG, buf)
	} else {
		return ret, errors.New("Touchscreen Type is not set or is unrecognized")
	}

	if err != nil {
		return ret, errors.New("Failed to read in the touch registers")
	}

	// Mask off the MSB and stick it to the LSB
	ret.X = int(buf[0]) & 0xf
	ret.X = ret.X << 8
	ret.X = ret.X | int(buf[1])

	ret.Y = int(buf[2] & 0xf)
	ret.Y = ret.Y << 8
	ret.Y = ret.Y | int(buf[3])

	if t.Debug {
		fmt.Printf("Touch Recevied: %d, %d\n", ret.X, ret.Y)
	}

	return ret, nil
}
