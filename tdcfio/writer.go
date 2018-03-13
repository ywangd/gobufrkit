package tdcfio

// Writer is the interface that wraps basic operations to serialise to output.
type Writer interface {
    WriteUint(v uint, n int) error
    WriteInt(v int, n int) error
    WriteBool(v bool) error

    // WriteBytes writes n bytes from an array of byte
    // The array is truncated if it has more than n elements. If the array
    // is shorter, null byte is used till n bytes is written.
    WriteBytes(v []byte, n int) error
    WriteBinary(v *Binary, n int) error

    // Write a float64 in format of IEE754 float32
    WriteFloat32(v float64) error
}