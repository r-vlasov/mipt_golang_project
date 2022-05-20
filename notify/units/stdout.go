package notifyunit

import (
	"fmt"
)

type StdoutNotifier struct {
}

func (n *StdoutNotifier) Notify(msg string) {
	fmt.Println(msg)
}
