// package tdcfio provides basic io features for working with WMO TDCF messages.
package tdcfio

const NBITS_PER_BYTE = 8

type InputType int

const (
    BinaryInput   InputType = iota
    FlatTextInput
    FlatJsonInput
)

// Reader is the interface that wraps the basic operations needed for deserialize a value.
//
// Some ReadXXX operations take an integer argument to specify the width of the
// returned value. This width could be of any unit, including bit, byte or token.
// The implementation could also choose to ignore this argument if it is not
// necessary in corresponding situation.
type Reader interface {
    // Pos returns the position for next read
    Pos() int

    ReadNumber(n int) (float64, error)
    ReadUint(n int) (uint, error)
    ReadInt(n int) (int, error)
    ReadBool() (bool, error)
    ReadBytes(n int) ([]byte, error)
    ReadBinary(n int) (*Binary, error)

    // ReadFloat32 reads a IEE754 float32 value and returns it as a float64
    ReadFloat32() (float64, error)
}

// PeekableReader is an augmented tdcfio.Reader by providing peeking features.
//
// The peeking operations can peek a value by skipping number of units (NOT uint).
// The unit can be anything such as bit, bytes or tokens depending on the actual
// implementation.
type PeekableReader interface {
    Reader
    // PeekUint returns an uint value by skipping number of skip unit.
    PeekUint(skip int, n int) (uint, error)

    // PeekBytes returns an array of byte by skipping number of skip unit.
    PeekBytes(skip int, n int) ([]byte, error)
}
