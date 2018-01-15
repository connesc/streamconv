package streamconv

import "io"

type Splitter interface {
	ReadItem() (item io.Reader, err error)
}

type Converter interface {
	Convert(input io.Reader) (output io.Reader, err error)
}

type Joiner interface {
	WriteItem(item io.Reader) (err error)
}
