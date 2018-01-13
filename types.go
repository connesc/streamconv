package streamconv

type ItemReader interface {
	ReadItem() (item []byte, err error)
}

type Converter interface {
	Convert(src []byte) (dst []byte, err error)
}

type ItemWriter interface {
	WriteItem(item []byte) (err error)
}
