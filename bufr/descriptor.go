package bufr

type Descriptor interface {
    Name()
    Id() string
    F() string
    X() string
    Y() string
}

type BaseDescriptor struct {
    id string
}

func (d *BaseDescriptor) Id() string {
    return d.id
}

func (d *BaseDescriptor) F() string {
    return d.id[0:1]
}

func (d *BaseDescriptor) X() string {
    return d.id[1:3]
}

func (d *BaseDescriptor) Y() string {
    return d.id[3:]
}
