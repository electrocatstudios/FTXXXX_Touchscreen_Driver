# FTXXXX_Touchscreen_Driver
Go bindings for touchscreen implementations for FT62xx and FT5x06 touchscreens

## Usage

Get the library into your GOPATH
```
go get github.com/electrocatstudios/FTXXX_Touchscreen_Driver
```
Then in your file include the following:

```
import touchscreen "github.com/electrocatstudios/FTXXXX_Touchscreen_Driver"
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
