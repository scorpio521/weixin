package tools
import (
	"fmt"
)

//写type的时候第一个字母要大写
type DivideError struct {
	Code  uint32
	Msg   string
}
func (de *DivideError) Error() string {
	return fmt.Sprintf("code = %d ; msg = %s", de.Code, de.Msg)
}

