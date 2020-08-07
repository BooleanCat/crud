package crud

import (
	"io"
	"os"
)

func Create(name string, content io.Reader) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return err
}

func Update(name string, content io.Reader) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return err
}

func Delete(name string) error {
	if _, err := os.Stat(name); err != nil {
		return err
	}

	return os.RemoveAll(name)
}

func ReadInto(name string, w io.Writer) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	return err
}

type NotFoundErr struct{}

func (e NotFoundErr) Error() string {
	return "not found"
}

var _ error = NotFoundErr{}
