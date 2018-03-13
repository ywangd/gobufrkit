package payload

import (
    "github.com/pkg/errors"
    "github.com/ywangd/gobufrkit/bufr"
    "fmt"
)

const refNotSet = -1

// Bitmap wraps the information about a bitmap, i.e. its bits etc.
type Bitmap struct {
    // start and stop (exclusive) indices for the bitmap nodes (031031)
    Index0 int
    Index1 int
}

// BitmapManager provides operations to lifecycle management of bitmaps
// and different type of bitmaps (ad-hoc and reuse)
type BitmapManager struct {
    cellsBuilder CellsBuilder
    // End index (inclusive) of the target descriptors
    backref1    int
    targetNodes []*bufr.ValuedNode

    reusableBtimap *Bitmap
    currentBitmap  *Bitmap
}

// NewAssessment signifies the beginning of a new assessment session
// This may lead to backref1 to be set.
func (bm *BitmapManager) NewAssessment() {
    if bm.backref1 == refNotSet {
        bm.backref1 = bm.cellsBuilder.Len() - 1
    }
}

// NewBitmap signifies the beginning of a new bitmap definition session
func (bm *BitmapManager) NewBitmap(reusable bool) {
    bitmap := &Bitmap{Index0: bm.cellsBuilder.Len()}
    bm.currentBitmap = bitmap
    if reusable {
        bm.reusableBtimap = bitmap
    }
}

// EndBitmap signifies the end of a bitmap definition session
func (bm *BitmapManager) EndBitmap() {
    bm.currentBitmap.Index1 = bm.cellsBuilder.Len()
}

func (bm *BitmapManager) RecallBitmap() {
    bm.currentBitmap = bm.reusableBtimap
}

// CancelBitmap cancels the definition of reusable bitmap 237255
func (bm *BitmapManager) CancelBitmap() {
    bm.reusableBtimap = nil
    bm.currentBitmap = nil
}

// CancelBackref cancels all backward reference 235000
func (bm *BitmapManager) CancelBackref() {
    bm.backref1 = refNotSet
    bm.CancelBitmap()
}

// Bits returns the bit values of the current bitmap
func (bm *BitmapManager) Bits() ([]uint, error) {
    i0, i1 := bm.currentBitmap.Index0, bm.currentBitmap.Index1
    bits := make([]uint, i1-i0)
    for i := i0; i < i1; i++ {
        cell := bm.cellsBuilder.Cell(i)
        b, err := cell.UintValue()
        if err != nil {
            return nil, errors.Wrap(err, "bitmap bit is not an unit")
        }
        bits[i-i0] = b
    }
    return bits, nil
}

// TODO: cache for optimisation
func (bm *BitmapManager) InitTargetNodes() ([]*bufr.ValuedNode, error) {
    bits, err := bm.Bits()
    if err != nil {
        return nil, err
    }
    var nodes []*bufr.ValuedNode
    for i := bm.backref1; i >= 0; i-- {
        nodes = append([]*bufr.ValuedNode{bm.cellsBuilder.Cell(i).Node()}, nodes...)
        if len(nodes) < len(bits) {
            continue
        }
        break
    }

    if len(bits) != len(nodes) {
        return nil, fmt.Errorf(
            "inconsistent number of bits and target nodes: %v != %v",
            len(bits), len(nodes))
    }
    bm.targetNodes = []*bufr.ValuedNode{}
    for i, b := range bits {
        if b == 0 {
            bm.targetNodes = append(bm.targetNodes, nodes[i])
        }
    }
    return bm.targetNodes, nil
}

func (bm *BitmapManager) NextTargetNode() (node *bufr.ValuedNode) {
    node, bm.targetNodes = bm.targetNodes[0], bm.targetNodes[1:]
    return
}
