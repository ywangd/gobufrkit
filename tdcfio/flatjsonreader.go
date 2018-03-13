package tdcfio

import (
    "encoding/json"
    "github.com/pkg/errors"
    "fmt"
    "io"
)

// FlatJsonReader implements the tdcfio.Reader interface for reading from JSON messages.
// TODO: The width argument could be used to ensure the value in binary is within given width
type FlatJsonReader struct {
    d   *json.Decoder
    pos int
}

// NewFlatJsonReader returns a pointer to FlatJsonReader.
//
// JSON streams are read in unit of json.Token. Hence the width argument in some ReadXXX
// operations are not needed and ignored.
func NewFlatJsonReader(reader io.Reader) *FlatJsonReader {
    d := json.NewDecoder(reader)
    return &FlatJsonReader{d: d}
}

func (r *FlatJsonReader) Pos() int {
    return r.pos
}

func (r *FlatJsonReader) ReadNumber(n int) (float64, error) {
    t, err := r.token()
    if err != nil {
        return 0, err
    }
    r.pos += n
    return r.float64(t)
}

func (r *FlatJsonReader) ReadUint(n int) (uint, error) {
    t, err := r.token()
    if err != nil {
        return 0, err
    }
    r.pos += n
    return r.uint(t)
}

func (r *FlatJsonReader) ReadInt(n int) (int, error) {
    t, err := r.token()
    if err != nil {
        return 0, err
    }
    r.pos += n
    return r.int(t)
}

func (r *FlatJsonReader) ReadBool() (bool, error) {
    t, err := r.token()
    if err != nil {
        return false, err
    }
    r.pos += 1
    return r.bool(t)
}

func (r *FlatJsonReader) ReadBytes(n int) ([]byte, error) {
    t, err := r.token()
    if err != nil {
        return nil, err
    }
    r.pos += n * NBITS_PER_BYTE
    return r.bytes(t)
}

func (r *FlatJsonReader) ReadBinary(n int) (*Binary, error) {
    t, err := r.token()
    if err != nil {
        return nil, err
    }
    r.pos += n
    return r.binary(t)
}

func (r *FlatJsonReader) ReadFloat32() (float64, error) {
    t, err := r.token()
    if err != nil {
        return 0, err
    }
    r.pos += 32
    return r.float64(t)
}

// token gets next non-delimiter token from the decoder
func (r *FlatJsonReader) token() (json.Token, error) {
    for {
        t, err := r.d.Token()
        if err != nil {
            return nil, errors.Wrap(err, "cannot get token")
        }
        switch t.(type) {
        case json.Delim:
        default:
            return t, nil
        }
    }
}

// unit attempts to get an uint from the given token
func (r *FlatJsonReader) uint(t json.Token) (uint, error) {
    v, err := r.float64(t)
    if err != nil {
        return 0, err
    }
    return uint(v), nil
}

func (r *FlatJsonReader) int(t json.Token) (int, error) {
    v, err := r.float64(t)
    if err != nil {
        return 0, err
    }
    return int(v), nil
}

func (r *FlatJsonReader) bool(t json.Token) (bool, error) {
    v, ok := t.(bool)
    if !ok {
        return false, fmt.Errorf("value is not bool type: %v, %T", t, t)
    }
    return v, nil
}

func (r *FlatJsonReader) bytes(t json.Token) ([]byte, error) {
    s, err := r.string(t)
    if err != nil {
        return nil, err
    }
    return []byte(s), nil
}

func (r *FlatJsonReader) binary(t json.Token) (*Binary, error) {
    s, err := r.string(t)
    if err != nil {
        return nil, err
    }
    return NewBinaryFromString(s)
}

func (r *FlatJsonReader) string(t json.Token) (string, error) {
    v, ok := t.(string)
    if !ok {
        return "", fmt.Errorf("value is not string type: %v, %T", t, t)
    }
    return v, nil
}

func (r *FlatJsonReader) float64(t json.Token) (float64, error) {
    v, ok := t.(float64)
    if !ok {
        return 0, fmt.Errorf("value is not float64 type: %v, %T", t, t)
    }
    return v, nil
}

///////////////////////////////////////////////////////////////////////////////
type PeekableFlatJsonReader struct {
    *FlatJsonReader
    tokens []json.Token
}

// NewPeekableFlatJsonReader returns a pointer to PeekableFlatJsonReader.
//
// PeekableFlatJsonReader implements tdcfio.PeekableReader for reading from JSON streams.
func NewPeekableFlatJsonReader(reader io.Reader) *PeekableFlatJsonReader {
    return &PeekableFlatJsonReader{FlatJsonReader: NewFlatJsonReader(reader)}
}

func (pr *PeekableFlatJsonReader) ReadNumber(n int) (float64, error) {
    if t := pr.shift(); t != nil {
        pr.pos += n
        return pr.float64(t)
    }
    return pr.FlatJsonReader.ReadNumber(n)
}

func (pr *PeekableFlatJsonReader) ReadUint(n int) (uint, error) {
    if t := pr.shift(); t != nil {
        pr.pos += n
        return pr.uint(t)
    }
    return pr.FlatJsonReader.ReadUint(n)
}

func (pr *PeekableFlatJsonReader) ReadInt(n int) (int, error) {
    if t := pr.shift(); t != nil {
        pr.pos += n
        return pr.int(t)
    }
    return pr.FlatJsonReader.ReadInt(n)
}

func (pr *PeekableFlatJsonReader) ReadBool() (bool, error) {
    if t := pr.shift(); t != nil {
        pr.pos += 1
        return pr.bool(t)
    }
    return pr.FlatJsonReader.ReadBool()
}

func (pr *PeekableFlatJsonReader) ReadBytes(n int) ([]byte, error) {
    if t := pr.shift(); t != nil {
        pr.pos += n * NBITS_PER_BYTE
        return pr.bytes(t)
    }
    return pr.FlatJsonReader.ReadBytes(n)
}

func (pr *PeekableFlatJsonReader) ReadBinary(n int) (*Binary, error) {
    if t := pr.shift(); t != nil {
        pr.pos += n
        return pr.binary(t)
    }
    return pr.FlatJsonReader.ReadBinary(n)
}

func (pr *PeekableFlatJsonReader) ReadFloat32() (float64, error) {
    if t := pr.shift(); t != nil {
        pr.pos += 32
        return pr.float64(t)
    }
    return pr.ReadFloat32()
}

func (pr *PeekableFlatJsonReader) PeekUint(skip int, n int) (uint, error) {
    t, err := pr.peek(skip)
    if err != nil {
        return 0, err
    }
    return pr.uint(t)
}

func (pr *PeekableFlatJsonReader) PeekBytes(skip int, n int) ([]byte, error) {
    t, err := pr.peek(skip)
    if err != nil {
        return nil, err
    }
    return pr.bytes(t)
}

func (pr *PeekableFlatJsonReader) peek(skip int) (json.Token, error) {
    // number of tokens need to be peeked
    n := skip - len(pr.tokens) + 1
    for i := 0; i < n; i++ {
        t, err := pr.token()
        if err != nil {
            return nil, errors.Wrap(err, "cannot peek token")
        }
        pr.tokens = append(pr.tokens, t)
    }
    return pr.tokens[len(pr.tokens)-1], nil
}

// shift returns the first token from the tokens cache or nil if none is available
func (pr *PeekableFlatJsonReader) shift() (t json.Token) {
    if len(pr.tokens) == 0 {
        return nil
    }
    t, pr.tokens = pr.tokens[0], pr.tokens[1:]
    return t
}
