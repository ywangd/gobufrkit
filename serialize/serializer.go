package serialize

import (
    "github.com/ywangd/gobufrkit/bufr"
    "io"
)

type Serializer interface {
    Serialize(message *bufr.Message) error
}

type FlatTextSerializer struct {
    v *FlatTextVisitor
}

func NewFlatTextSerializer(writer io.Writer) *FlatTextSerializer {
    return &FlatTextSerializer{
        v: NewFlatTextVisitor(writer),
    }
}

func (s *FlatTextSerializer) Serialize(message *bufr.Message) error {
    return message.Accept(s.v)
}
