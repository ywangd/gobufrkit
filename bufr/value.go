package bufr

type Value interface {
    Descriptor() Descriptor
    Value() interface{}
    Attributes() []Value
}
