package tdcfio

import (
    "fmt"
    "strings"
    "encoding/json"
    "strconv"
)

// Binary represent a number of binary bits
//
// The bits are stored as a slice of byte from the leftmost digit to the rightmost
// i.e. the first byte represent the first 8 (leftmost) bits
type Binary struct {
    b     []byte
    nbits int
}

func NewBinary(b []byte, nbits int) (*Binary, error) {

    nbytes := nbits / NBITS_PER_BYTE
    if x := nbits % NBITS_PER_BYTE; x > 0 {
        nbytes += 1
    }

    if len(b) != nbytes {
        return nil, fmt.Errorf("number of bytes and bits mismatch")
    }

    return &Binary{b, nbits}, nil

}

func NewBinaryFromString(s string) (*Binary, error) {
    nbits := len(s)
    nbitsSaved := nbits
    b := []byte{}

    for nbits >= NBITS_PER_BYTE {
        u, err := strconv.ParseUint(s[:NBITS_PER_BYTE], 2, NBITS_PER_BYTE)
        if err != nil {
            return nil, err
        }
        b = append(b, byte(u))
        s = s[NBITS_PER_BYTE:]
        nbits -= NBITS_PER_BYTE
    }
    if nbits > 0 {
        u, err := strconv.ParseUint(s, 2, NBITS_PER_BYTE)
        if err != nil {
            return nil, err
        }
        b = append(b, byte(u))
    }

    return NewBinary(b, nbitsSaved)
}

func NewBinaryFromUint(v uint, nbits int) (*Binary, error) {
    bs := []byte{}
    for nbitsShift := uint(nbits - NBITS_PER_BYTE); nbitsShift >= 0; nbitsShift -= uint(NBITS_PER_BYTE) {
        bs = append(bs, byte(((127<<nbitsShift)&v)>>nbitsShift))
    }
    binary, err := NewBinary(bs, nbits)
    if err != nil {
        return nil, err
    }
    return binary, nil
}

func (bin *Binary) String() string {

    var ret []string
    var length int
    for i, b := range bin.b {
        length = NBITS_PER_BYTE
        if i == len(bin.b)-1 {
            if x := bin.nbits % NBITS_PER_BYTE; x != 0 {
                length = x
            }
        }
        ret = append(ret, fmt.Sprintf(fmt.Sprintf("%%0%db", length), b))
    }
    return strings.Join(ret, "")
}

// Get value at bit n (bit 0 is the most significant bit)
// TODO: return an error for out of bound indices?
func (bin *Binary) Bit(n int) bool {
    if n >= bin.nbits {
        return false
    }

    ibyte := n / NBITS_PER_BYTE
    ibit := n % NBITS_PER_BYTE

    b := bin.b[ibyte]
    padding := 0
    if ibyte == len(bin.b)-1 && bin.nbits%NBITS_PER_BYTE != 0 {
        padding = NBITS_PER_BYTE - (bin.nbits % NBITS_PER_BYTE)
    }

    return b>>uint(NBITS_PER_BYTE-ibit-1-padding)&1 == 1
}

func (bin *Binary) Nbits() int {
    return bin.nbits
}

func (bin *Binary) MarshalJSON() ([]byte, error) {
    return json.Marshal(bin.String())
}

func (bin *Binary) UnmarshalJSON(data []byte) error {
    s := ""
    json.Unmarshal(data, &s)

    b, err := NewBinaryFromString(s)
    if err != nil {
        return err
    }
    bin.nbits = b.nbits
    bin.b = b.b

    return nil
}
