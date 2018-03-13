package payload

import (
    "github.com/ywangd/gobufrkit/table"
    "math"
    "github.com/pkg/errors"
    "github.com/ywangd/gobufrkit/bufr"
    "fmt"
)

func calcPackingInfo(v *DesVisitor, descriptor table.Descriptor) (*bufr.PackingInfo, error) {
    entry := descriptor.Entry().(*table.Bentry)

    info := &bufr.PackingInfo{Unit: entry.Unit}

    switch entry.Unit {
    case table.STRING:
        // TODO: Assert scale and refval of entry are zero
        info.Scale = 0
        info.Refval = 0
        if v.nbitsString != 0 {
            info.Nbits = v.nbitsString
        } else {
            info.Nbits = entry.Nbits
        }

    case table.NUMERIC:
        info.Scale = entry.Scale + v.scaleOffset + v.scaleIncrement

        newRefvalNode, ok := v.newRefvalNodes[descriptor.Id()]
        if ok {
            if !v.cellsBuilder.LastCellEquality() {
                return nil, fmt.Errorf("new ref val not equal across all (compressed) subsets")
            }
            newRefval, err := v.cellsBuilder.Cell(newRefvalNode.Index).FloatValue()
            if err != nil {
                return nil, errors.Wrap(err, "cannot get new ref val")
            }
            info.Refval = newRefval * math.Pow10(v.refvalFactor)
        } else {
            info.Refval = float64(entry.Refval) * math.Pow10(v.refvalFactor)
        }
        info.Nbits = entry.Nbits + v.nbitsOffset + v.nbitsIncrement

    default:
        info.Scale = entry.Scale
        info.Refval = float64(entry.Refval)
        info.Nbits = entry.Nbits
    }

    return info, nil
}
