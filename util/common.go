package util

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"github.com/mattn/go-isatty"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func Zip(writer io.Writer, reader io.Reader) error {
	zw := gzip.NewWriter(writer)
	defer zw.Close()

	for buf, r := make([]byte, 65536), bufio.NewReader(reader); ; {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = zw.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return zw.Flush()
}

func GenerateTraceID() string {
	return uuid.NewV4().String()
}

func CheckIfTerminal(out io.Writer) bool {
	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		return false
	}
	return true
}

func JsonString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
