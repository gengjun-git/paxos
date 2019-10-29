package log

import "fmt"

func Log(format string, a ...interface{})  {
	fmt.Println(fmt.Sprintf(format, a...))
}
