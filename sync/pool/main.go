package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

//TODO: create pool of bytes.Buffers which can be reused.

var bufferPool = sync.Pool{
	New: func() interface{} {
		fmt.Println("Creating new buffer")
		return new(bytes.Buffer)
	},
}

func log(w io.Writer, val string) {
	b := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(b)
	b.Reset()
	b.WriteString(time.Now().Format("15:04:05"))
	b.WriteString(" : ")
	b.WriteString(val)
	b.WriteString("\n")
	w.Write(b.Bytes())
}

func main() {
	log(os.Stdout, "debug-string1")
	log(os.Stdout, "debug-string2")
}
