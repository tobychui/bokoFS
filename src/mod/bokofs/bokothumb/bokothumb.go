package bokothumb

import (
	"errors"
	"io"
	"net/http"
	"os"
)

type File struct {
	http.File
	io.Writer
}

func (f *File) Write(p []byte) (n int, err error) {
	return 0, errors.New("readonly file system")
}

func (f *File) Close() error {
	return f.File.Close()
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.File.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.File.Seek(offset, whence)
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return f.File.Readdir(count)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.File.Stat()
}
