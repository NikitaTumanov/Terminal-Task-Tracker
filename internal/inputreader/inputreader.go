package inputreader

import (
	"bufio"
	"os"
)

type Reader struct {
	Input string
}

func (r *Reader) Read() {
	reader := bufio.NewReader(os.Stdin)
	r.Input, _ = reader.ReadString('\n')
}
