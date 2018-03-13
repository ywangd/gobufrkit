//go:generate stringer -type=Operator
package table

import (
    "fmt"
)

const (
    F_ELEMENT     = iota
    F_REPLICATION
    F_OPERATOR
    F_SEQUENCE
)

// Operator descriptor types
type Operator int

const (
    OP_NBITS_OFFSET              Operator = 201
    OP_SCALE_OFFSET              Operator = 202
    OP_NEW_REFVAL                Operator = 203
    OP_ASSOCIATE_FIELD           Operator = 204
    OP_INSERT_STRING             Operator = 205
    OP_SKIP_LOCAL                Operator = 206
    OP_MODIFY_PACKING            Operator = 207
    OP_SET_STRING_LENGTH         Operator = 208
    OP_DATA_NOT_PRESENT          Operator = 221
    OP_QUALITY_INFO              Operator = 222
    OP_SUBSTITUTION              Operator = 223
    OP_FIRST_ORDER_STATS         Operator = 224
    OP_DIFFERENCE_STATS          Operator = 225
    OP_REPLACEMENT               Operator = 232
    OP_CANCEL_BACK_REFERENCE     Operator = 235
    OP_DEFINE_BITMAP             Operator = 236
    OP_RECALL_BITMAP             Operator = 237
    OP_DEFINE_EVENT              Operator = 241
    OP_DEFINE_CONDITIONING_EVENT Operator = 242
    OP_CATEGORICAL_VALUES        Operator = 243
)

// ID represents the very basic information that a descriptor can provide,
// i.e. it's integer type ID and the F X Y semantics. It does NOT get any
// further into the specific meanings of each descriptors.
type ID int

const (
    ID_031021 ID = 31021 // Associated field significance
    ID_031031 ID = 31031 // Bitmap data present indicator
    ID_008023 ID = 8023  // First order stats significance
    ID_008024 ID = 8024  // Difference stats significance

    ID_203255 ID = 203255
    ID_222000 ID = 222000
    ID_223000 ID = 223000
    ID_223255 ID = 223255
    ID_224000 ID = 224000
    ID_224255 ID = 224255
    ID_225000 ID = 225000
    ID_225255 ID = 225255
    ID_232000 ID = 232000
    ID_232255 ID = 232255
    ID_236000 ID = 236000
    ID_237000 ID = 237000
    ID_237255 ID = 237255
)
const CLASS_33_QA_INFO = 33

func (id ID) String() string {
    return fmt.Sprintf("%06d", id)
}

func (id ID) F() int {
    return int(id) / 100000
}

func (id ID) X() int {
    return (int(id) / 1000) % 100
}

func (id ID) Y() int {
    return int(id) % 1000
}

// The Descriptor interface is a general form to represent a BUFR descriptor.
// In addition to the F X Y semantics, it also has an associated Entry
// whose actual type depends on the descriptor type.
type Descriptor interface {
    Id() ID
    F() int
    X() int
    Y() int

    Operator() Operator
    Operand() int

    // Additional information about the descriptor which is mostly from
    // an entry in corresponding BUFR table.
    Entry() Entry
}

type BaseDescriptor struct {
    ID
    entry Entry
}

func (d *BaseDescriptor) String() string {
    if d.entry != nil {
        return fmt.Sprintf("%v %v", d.Id(), d.entry.Name())
    }
    return d.ID.String()
}

func (d *BaseDescriptor) Id() ID {
    return d.ID
}

func (d *BaseDescriptor) Operator() Operator {
    return Operator(int(d.ID) / 1000)
}

func (d *BaseDescriptor) Operand() int {
    return d.Y()
}

func (d *BaseDescriptor) Entry() Entry {
    return d.entry
}

type ElementDescriptor struct {
    BaseDescriptor
}

type ReplicationDescriptor struct {
    BaseDescriptor
}

type OperatorDescriptor struct {
    BaseDescriptor
}

type SequenceDescriptor struct {
    BaseDescriptor
}

func NewElementDescriptor(id ID, entry Entry) *ElementDescriptor {
    return &ElementDescriptor{BaseDescriptor{id, entry}}
}

func NewSequenceDescriptor(id ID, entry Entry) *SequenceDescriptor {
    return &SequenceDescriptor{BaseDescriptor{id, entry}}
}

// RootDescriptor is a singleton placeholder root descriptor
var RootDescriptor = &BaseDescriptor{ID: 0, entry: &Rentry{name: "Root"}}

// NewLocalDescriptor returns a placeholder element descriptor for unknown local descriptor.
// An unknown local descriptor can appear immediately after operator 206YYY
func NewLocalDescriptor(id ID) *ElementDescriptor {
    return NewElementDescriptor(id, &Bentry{name: "LOCAL DESCRIPTOR"})
}

// TODO: Replace LOCAL descriptor with generic decorate descriptor
type DecorateDescriptor struct {
    Descriptor Descriptor
    Initial    int
    Name       string
}

func (dd *DecorateDescriptor) String() string {
    return fmt.Sprintf("%s%s %v", string(dd.Initial), dd.Id().String()[1:], dd.Name)
}

func (dd *DecorateDescriptor) Id() ID {
    return dd.Descriptor.Id()
}

func (dd *DecorateDescriptor) F() int {
    return dd.Descriptor.F()
}

func (dd *DecorateDescriptor) X() int {
    return dd.Descriptor.X()
}

func (dd *DecorateDescriptor) Y() int {
    return dd.Descriptor.Y()
}

func (dd *DecorateDescriptor) Operator() Operator {
    return dd.Descriptor.Operator()
}

func (dd *DecorateDescriptor) Operand() int {
    return dd.Descriptor.Operand()
}

func (dd *DecorateDescriptor) Entry() Entry {
    return dd.Descriptor.Entry()
}
