package streamconv

import "io"

type ItemReader interface {
	ReadItem() (item io.Reader, err error)
}

type Converter interface {
	Convert(src io.Reader) (dst io.Reader, err error)
}

type Transformer interface {
	Transform(src ItemReader) (dst ItemReader, err error)
}

type ItemWriter interface {
	WriteItem(item io.Reader) (err error)
}

func Copy(dst ItemWriter, src ItemReader) (err error) {
	for {
		item, err := src.ReadItem()
		if err != nil {
			return err
		}

		err = dst.WriteItem(item)
		if err != nil {
			return err
		}
	}
}
