# FTXXXX_Touchscreen_Driver
Golang bindings for touchscreen implementations for FT62xx and FT5x06 touchscreens

## Usage

Get the library into your GOPATH
```
go get github.com/electrocatstudios/FTXXX_Touchscreen_Driver
```
Then in your file include the following:

```
import "github.com/electrocatstudios/FTXXXX_Touchscreen_Driver/touchscreen"
```

## Code Example

```
t := touchscreen.Touchscreen{}
t.Init()

numTouches, _ := t.GetTouchesCount()
if numTouches > 0 {
    touch, _ := t.GetTouches()
    fmt.Printf("X: %d, Y: %d\n", touch.X, touch.Y)
}
```

## Note on coordinates
It may be necessary to translate the points received. The position of the points is presented in portrait mode so you may need to rotate depending on your screen orientation. On the ili9341 screen in landscape mode (for example) the points are as follows:

(0,320)                   (0,0)
X-------------------------X
|                         |
|                         |
|                         |
|                         |
|                         |
X-------------------------X
(240,320)                 (240,0)