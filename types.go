package streamconv

import "io"

type ItemReader interface {
	ReadItem() (item io.Reader, err error)
}

type Converter interface {
	Convert(src io.Reader) (dst io.Reader, err error)
}

type ItemWriter interface {
	WriteItem(item io.Reader) (err error)
}
