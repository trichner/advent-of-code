package term

import (
	"fmt"
	"io"
)

type Color string

var (
	ColorReset   Color = "\033[0m"
	ColorRed     Color = "\033[31m"
	ColorGreen   Color = "\033[32m"
	ColorYellow  Color = "\033[33m"
	ColorBlue    Color = "\033[34m"
	ColorMagenta Color = "\033[35m"
	ColorCyan    Color = "\033[36m"
	ColorGray    Color = "\033[37m"
	ColorWhite   Color = "\033[97m"
)

func StringInColor(s string, c Color) string {
	return fmt.Sprintf("%s%s%s", c, s, ColorReset)
}

func WriteInColor(w io.Writer, s string, c Color) error {
	_, err := fmt.Fprintf(w, "%s%s%s", c, s, ColorReset)
	return err
}
