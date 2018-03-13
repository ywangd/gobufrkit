package bufr

import (
    "fmt"
    "encoding/json"
)

// LookupFunction is for looking up meanings of code values on code tables
type LookupFunction = func(uint) string

var MISSING_VALUE_OF_NBITS = make(map[int]uint)

func init() {
    // Store the missing values for corresponding number of bits for
    // quick lookup.
    for i := 1; i <= 64; i++ {
        MISSING_VALUE_OF_NBITS[i] = 1<<uint(i) - 1
    }
}

// IsMissing checks whether the given value is a missing value for the given number of
// bits. ONLY uint values can be missing. This is NOT checking whether the value is nil.
func IsMissing(value interface{}, nbits int) bool {
    if _, ok := value.(uint); ok && nbits > 1 {
        return value == MISSING_VALUE_OF_NBITS[nbits]
    }
    return false
}

// MissingValue returns uint number equals to the missing value of given number of bits
func MissingValue(nbits int) (uint, error) {
    v, ok := MISSING_VALUE_OF_NBITS[nbits]
    if ok {
        return v, nil
    }
    return 0, fmt.Errorf("invalid number of bits for retrieving missing value: %v", nbits)
}

// Field represents a field of a BUFR message's section. Its value can be anything,
// including primitive types and custom types for UnexpandedTemplate and Payload.
type Field struct {
    // Name of the field
    Name string

    // Value of the field. The Value can be a custom type Payload which covers
    // the entire payload (the meat of the data section) of a BUFR message.
    Value interface{}

    // Number of nbits taken by this field.
    Nbits int

    // Lookup function for code values
    Lookup LookupFunction

    // TODO: Maybe the Field can be an interface and have Hidden and Virtual as implementations?
    // A hidden field does not show up in certain serialised format, e.g. Text.
    Hidden bool

    // A virtual field is something that does not actually have corresponding
    // bits in the source file, e.g. a computed field
    Virtual bool
}

func (f *Field) Accept(visitor Visitor) error {
    return visitor.VisitField(f)
}

func (f *Field) IsMissing() bool {
    return IsMissing(f.Value, f.Nbits)
}

func (f *Field) MarshalJSON() ([]byte, error) {
    return json.Marshal(f.Value)
}

func NewField(name string, value interface{}, nbits int) *Field {
    return &Field{Name: name, Value: value, Nbits: nbits, Virtual: false}
}

func NewHiddenField(name string, value interface{}, nbits int) *Field {
    field := NewField(name, value, nbits)
    field.Hidden = true
    return field
}
