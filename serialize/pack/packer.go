package pack

import (
    "github.com/ywangd/gobufrkit/bufr"
    "github.com/ywangd/gobufrkit/table"
    "github.com/ywangd/gobufrkit/tdcfio"
    "fmt"
    "math"
)

type Packer interface {
    Pack(*bufr.ValuedNode, interface{}) error
}

type UncompressedPacker struct {
    w tdcfio.Writer
}

func NewUncompressedPacker(w tdcfio.Writer) Packer {
    return &UncompressedPacker{w: w}
}

func (p *UncompressedPacker) Pack(node *bufr.ValuedNode, value interface{}) error {
    var (
        err error
        ok  bool
    )
    info := node.PackingInfo

    switch info.Unit {
    case table.STRING:
        v, ok := value.([]byte)
        if !ok {
            return fmt.Errorf("value is not a string: %v", value)
        }
        return p.w.WriteBytes(v, info.Nbits/tdcfio.NBITS_PER_BYTE)

    case table.NONNEG_CODE, table.FLAG:
        var v uint
        if value == nil {
            v, err = bufr.MissingValue(info.Nbits)
            if err != nil {
                return err
            }
        } else {
            v, ok = value.(uint)
            if !ok {
                return fmt.Errorf("value is not an uint: %v", value)
            }
        }
        return p.w.WriteUint(v, info.Nbits)

    case table.NUMERIC:
        if value == nil {
            v, err := bufr.MissingValue(info.Nbits)
            if err != nil {
                return err
            }
            return p.w.WriteUint(v, info.Nbits)
        } else {
            switch value.(type) {
            case int:
                if info.Refval != 0 || info.Scale != 0 {
                    return fmt.Errorf(
                        "integer numeric value must have zero for both refval (%v) and scale (%v)",
                        info.Refval, info.Scale)
                }
                return p.w.WriteInt(value.(int), info.Nbits)

            case float64:
                xfloat := value.(float64)
                if info.Scale != 0 {
                    xfloat *= math.Pow10(info.Scale)
                }
                if info.Refval != 0 {
                    xfloat -= info.Refval
                }
                return p.w.WriteUint(uint(xfloat), info.Nbits)
            default:
                return fmt.Errorf("invalid data type for numeric value: %v", value)
            }
        }

    case table.BINARY:
        v, ok := value.(*tdcfio.Binary)
        if !ok {
            return fmt.Errorf("value is not of BINARY type: %v", value)
        }
        return p.w.WriteBinary(v, info.Nbits)

    default:
        return fmt.Errorf("unrecognised data unit: %v", info.Unit)
    }

    return nil
}

const NBITS_FOR_NBITS_DIFF = 6

type CompressedPacker struct {
    w tdcfio.Writer
}

func NewCompressedPacker(w tdcfio.Writer) Packer {
    return &CompressedPacker{w: w}
}

func (p *CompressedPacker) Pack(node *bufr.ValuedNode, values interface{}) error {

    info := node.PackingInfo
    // TODO: type assertion checking
    switch info.Unit {
    case table.STRING:
        p.w.WriteBytes(node.MinValue.([]byte), info.Nbits/tdcfio.NBITS_PER_BYTE)
        p.w.WriteInt(node.NbitsDiff, NBITS_FOR_NBITS_DIFF)
        // NbitsDiff is actually nbytes for string type
        if node.NbitsDiff != 0 {
            for _, v := range values.([]interface{}) {
                p.w.WriteBytes(v.([]byte), node.NbitsDiff)
            }
        }

    case table.NONNEG_CODE, table.FLAG:
        p.w.WriteUint(node.MinValue.(uint), info.Nbits)
        p.w.WriteInt(node.NbitsDiff, NBITS_FOR_NBITS_DIFF)
        if node.NbitsDiff != 0 {
            for _, v := range values.([]interface{}) {
                p.w.WriteUint(v.(uint)-node.MinValue.(uint), node.NbitsDiff)
            }
        }

    case table.NUMERIC:
        p.w.WriteUint(node.MinValue.(uint), info.Nbits)
        p.w.WriteInt(node.NbitsDiff, NBITS_FOR_NBITS_DIFF)
        if node.NbitsDiff != 0 {
            for _, v := range values.([]interface{}) {
                switch {
                case v == nil:
                    xuint, _ := bufr.MissingValue(node.NbitsDiff)
                    p.w.WriteUint(xuint, node.NbitsDiff)

                case info.Refval != 0 || info.Scale != 0:
                    xfloat := v.(float64)
                    if info.Scale != 0 {
                        xfloat *= math.Pow10(info.Scale)
                    }
                    if info.Refval != 0 {
                        xfloat -= info.Refval
                    }
                    p.w.WriteUint(uint(xfloat)-node.MinValue.(uint), node.NbitsDiff)

                default:
                    xint := v.(int)
                    p.w.WriteUint(uint(xint)-node.MinValue.(uint), node.NbitsDiff)
                }
            }
        }

    case table.BINARY:
        // TODO:
        panic("NYI")

    default:
        return fmt.Errorf("unrecognised data unit: %v", info.Unit)
    }

    return nil
}
