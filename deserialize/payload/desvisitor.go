package payload

import (
    "github.com/ywangd/gobufrkit/deserialize/unpack"
    "github.com/ywangd/gobufrkit/deserialize/ast"
    "github.com/ywangd/gobufrkit/tdcfio"
    "github.com/ywangd/gobufrkit/bufr"
    "github.com/pkg/errors"
    "github.com/ywangd/gobufrkit/table"
    "fmt"
)

type DesVisitorConfig struct {
    Compressed bool
    InputType  tdcfio.InputType
    // compatible mode for BUFRDC (i.e. insert 0 for some operator descriptors, e.g. 236000
    Compatible bool
    Verbose    bool
}

// DesVisitor is an implementation of ast.Visitor that constructs bufr.Subset by
// deserializing from an tdcfio.Reader.
type DesVisitor struct {
    config *DesVisitorConfig

    nsubsets int

    unpacker unpack.Unpacker

    treeBuilder  *TreeBuilder
    cellsBuilder CellsBuilder

    nbitsOffset    int                           // 201YYY
    scaleOffset    int                           // 202YYY
    newRefvalNodes map[table.ID]*bufr.ValuedNode // 203YYY
    assocPairs     *assocPairs                   // 204YYY
    nbitsIncrement int                           // 207YYY
    scaleIncrement int                           // 207YYY
    refvalFactor   int                           // 207YYY
    nbitsString    int                           // 208YYY

    bitmapManager *BitmapManager
}

func NewDeserializeVisitor(config *DesVisitorConfig, reader tdcfio.Reader, nsubsets int) (*DesVisitor, error) {
    unpacker, err := unpack.NewUnpacker(reader, nsubsets, config.Compressed, config.InputType)
    if err != nil {
        return nil, err
    }
    return &DesVisitor{
        config:   config,
        nsubsets: nsubsets,
        unpacker: unpacker,
    }, nil

}

// Produce creates and adds subsets to the given Payload
func (v *DesVisitor) Produce(payload *bufr.Payload) error {
    root, err := v.treeBuilder.Root()
    if err != nil {
        return err
    }
    v.cellsBuilder.Produce(payload, root)
    return nil
}

func (v *DesVisitor) reset() {
    v.treeBuilder = NewTreeBuilder(bufr.NewBlock())
    v.treeBuilder.Verbose = v.config.Verbose
    if v.config.Compressed && v.config.InputType == tdcfio.BinaryInput {
        v.cellsBuilder = &CompressedCellsBuilder{
            nsubsets: v.nsubsets,
            tissues:  make([]Tissue, v.nsubsets),
        }
    } else {
        v.cellsBuilder = &UncompressCellsBuilder{}
    }
    v.bitmapManager = &BitmapManager{cellsBuilder: v.cellsBuilder, backref1: refNotSet}

    v.nbitsOffset = 0
    v.scaleOffset = 0
    v.newRefvalNodes = nil
    v.assocPairs = &assocPairs{}
    v.nbitsIncrement = 0
    v.scaleIncrement = 0
    v.refvalFactor = 0
    v.nbitsString = 0
}

func (v *DesVisitor) VisitNode(node ast.Node) error {
    v.reset()
    for _, n := range node.Members() {
        if err := n.Accept(v); err != nil {
            return errors.Wrap(err, "cannot process member node")
        }
    }
    return nil
}

func (v *DesVisitor) VisitElementNode(node *ast.ElementNode) error {
    if node.NotPresent { // not processing associated field for not present data
        v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
        return nil
    }

    // associated fields (if any) must be processed first as their data appear first
    assocNodes, err := buildAssocNodes(v, node.Descriptor())
    if err != nil {
        return errors.Wrap(err, "cannot process associated fields")
    }

    enode, err := buildValuedNode(v, node.Descriptor())
    if err != nil {
        return errors.Wrap(err, "cannot process element descriptor")
    }
    if len(assocNodes) > 0 {
        for _, n := range assocNodes {
            enode.AddMember(n)
        }
    }
    return nil
}

// Descriptor 031021 (associated field significance) cannot have associated fields.
func (v *DesVisitor) VisitE031021Node(node *ast.E031021Node) error {
    enode, err := buildValuedNode(v, node.Descriptor())
    if err != nil {
        return errors.Wrap(err, "cannot process element descriptor 031021")
    }
    v.assocPairs.SetNode(enode)
    return nil
}

func (v *DesVisitor) VisitFixedReplicationNode(node *ast.FixedReplicationNode) error {
    v.treeBuilder.Push(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    defer v.treeBuilder.Pop()
    for i := 0; i < node.Descriptor().Y(); i++ {
        if err := buildBlock(v, node.Members()); err != nil {
            return errors.Wrap(err, "cannot process fixed replication node")
        }
    }
    return nil
}

func (v *DesVisitor) VisitDelayedReplicationNode(node *ast.DelayedReplicationNode) error {
    v.treeBuilder.Push(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    defer v.treeBuilder.Pop()
    if err := node.Members()[0].Accept(v); err != nil {
        return errors.Wrap(err, "cannot process delayed replication factor")
    }
    if !v.cellsBuilder.LastCellEquality() {
        return fmt.Errorf("delayed replication factor not equal across all (compressed) subsets")
    }
    cell := v.cellsBuilder.Cell(v.cellsBuilder.Len() - 1)
    nreplications, err := cell.UintValue()
    if err != nil {
        return errors.Wrap(err, "cannot get delayed replication factor value as uint")
    }
    for i := uint(0); i < nreplications; i++ { // loop of replication
        if err := buildBlock(v, node.Members()[1:]); err != nil {
            return errors.Wrap(err, "cannot process delayed replication")
        }
    }
    return nil
}

func (v *DesVisitor) VisitSequenceNode(node *ast.SequenceNode) error {
    v.treeBuilder.Push(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    defer v.treeBuilder.Pop()
    for _, m := range node.Members() {
        if err := m.Accept(v); err != nil {
            return errors.Wrap(err, "cannot process sequence members")
        }
    }
    return nil
}

func (v *DesVisitor) VisitOpNbitsOffsetNode(node *ast.OpNbitsOffsetNode) error {
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    operand := node.Descriptor().Operand()
    if operand == 0 {
        v.nbitsOffset = operand
    } else {
        v.nbitsOffset = operand - 128
    }
    return nil
}

func (v *DesVisitor) VisitOpScaleOffsetNode(node *ast.OpScaleOffsetNode) error {
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    operand := node.Descriptor().Operand()
    if operand == 0 {
        v.scaleOffset = operand
    } else {
        v.scaleOffset = operand - 128
    }
    return nil
}

func (v *DesVisitor) VisitOpNewRefvalNode(node *ast.OpNewRefvalNode) error {
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    operand := node.Descriptor().Operand()
    if operand != 255 {
        info := &bufr.PackingInfo{Unit: table.CODE, Nbits: operand}
        for _, m := range node.Members() {
            mnode, err := buildValuedNodeWithInfo(v, m.Descriptor(), info)
            if err != nil {
                return errors.Wrap(err, "cannot process new refval definition")
            }
            v.newRefvalNodes[m.Descriptor().Id()] = mnode
        }
    }
    // nothing to do for 203255 other than add the node to tree
    return nil
}

func (v *DesVisitor) VisitOpAssocFieldNode(node *ast.OpAssocFieldNode) error {
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    if node.Descriptor().Y() == 0 {
        v.assocPairs.Pop()
    } else {
        v.assocPairs.Push(node.Descriptor().Y())
    }
    for _, m := range node.Members() {
        if err := m.Accept(v); err != nil {
            return errors.Wrap(err, "cannot process assoc field significance")
        }
    }
    return nil
}

func (v *DesVisitor) VisitOpInsertStringNode(node *ast.OpInsertStringNode) error {
    // Decide this should not process associated fields
    info := &bufr.PackingInfo{Unit: table.STRING, Nbits: node.Descriptor().Y() * tdcfio.NBITS_PER_BYTE}
    _, err := buildValuedNodeWithInfo(v, node.Descriptor(), info)
    return err
}

func (v *DesVisitor) VisitOpSkipLocalNode(node *ast.OpSkipLocalNode) error {
    // Decide this should not process associated fields
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    info := &bufr.PackingInfo{Unit: table.BINARY, Nbits: node.Descriptor().Y()}
    _, err := buildValuedNodeWithInfo(v, node.Members()[0].Descriptor(), info)
    return err
}

func (v *DesVisitor) VisitOpModifyPackingNode(node *ast.OpModifyPackingNode) error {
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    operand := node.Descriptor().Operand()
    if operand == 0 {
        v.nbitsIncrement = 0
        v.scaleIncrement = 0
        v.refvalFactor = 0

    } else {
        v.nbitsIncrement = (10*operand + 2) / 3
        v.scaleIncrement = operand
        v.refvalFactor = operand
    }
    return nil
}

func (v *DesVisitor) VisitOpSetStringLengthNode(node *ast.OpSetStringLengthNode) error {
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    v.nbitsString = node.Descriptor().Y() * tdcfio.NBITS_PER_BYTE
    return nil
}

func (v *DesVisitor) VisitOpDataNotPresentNode(node *ast.OpDataNotPresentNode) error {
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: node.Descriptor()})
    for _, n := range node.Members() {
        if err := n.Accept(v); err != nil {
            return errors.Wrap(err, "cannot process data not presented node")
        }
    }
    return nil
}

func (v *DesVisitor) VisitOpAssessmentNode(node *ast.OpAssessmentNode) error {
    v.bitmapManager.NewAssessment()
    // insert a zero for this operator descriptor if in compatible mode
    buildZeroNode(v, node.Descriptor())
    // Construct bitmap
    if err := node.Bitmap.Accept(v); err != nil {
        return err
    }

    // Process all sandwiched attributes
    // In case some of the attribute descriptor are non-element descriptor
    // Find all attribute nodes using start and stop indices
    anodes, err := v.visitNodesForValuedNodes(node.Attrs)
    if err != nil {
        return errors.Wrap(err, "cannot deserialize sandwiched nodes")
    }

    // Initialise all target nodes
    targetNodes, err := v.bitmapManager.InitTargetNodes()
    if err != nil {
        return errors.Wrap(err, "cannot get bitmapping target nodes")
    }

    // Process bitmap source nodes, e.g. marker nodes
    snodes, err := v.visitNodesForValuedNodes(node.Members())
    if err != nil {
        return errors.Wrap(err, "cannot deserialize bitmapping source nodes")
    }
    for i, snode := range snodes {
        targetNodes[i].AddMember(snode)
        for _, anode := range anodes {
            snode.AddMember(anode)
        }
    }
    return nil
}

func (v *DesVisitor) VisitOpMarkerNode(node *ast.OpMarkerNode) error {
    targetNode := v.bitmapManager.NextTargetNode()
    fmt.Println("target node ", targetNode.Descriptor)
    packingInfo, err := calcPackingInfo(v, targetNode.Descriptor)
    if err != nil {
        return err
    }
    if targetNode.Descriptor.Operator() == table.OP_DIFFERENCE_STATS {
        refval := -(1 << uint(packingInfo.Nbits))
        packingInfo.Refval = float64(refval)
        packingInfo.Nbits += 1
    }
    _, err = buildValuedNodeWithInfo(v, node.Descriptor(), packingInfo)
    return err
}

func (v *DesVisitor) VisitOpCancelBackRefNode(node *ast.OpCancelBackRefNode) error {
    v.bitmapManager.CancelBackref()
    return nil
}

func (v *DesVisitor) VisitOpCancelBitmapNode(node *ast.OpCancelBitmapNode) error {
    v.bitmapManager.CancelBitmap()
    return nil
}

func (v *DesVisitor) VisitBitmapNode(node *ast.BitmapNode) error {
    if node.Descriptor() != nil {
        buildZeroNode(v, node.Descriptor())
    }
    if node.Descriptor().Id() == table.ID_237000 {
        v.bitmapManager.RecallBitmap()
        return nil
    }
    v.bitmapManager.NewBitmap(node.Descriptor() != nil)
    defer v.bitmapManager.EndBitmap()
    for _, m := range node.Members() {
        if err := m.Accept(v); err != nil {
            return err
        }
    }
    return nil
}

func (v *DesVisitor) visitNodesForValuedNodes(astNodes []ast.Node) ([]*bufr.ValuedNode, error) {
    a0 := v.cellsBuilder.Len()
    for _, attr := range astNodes {
        if err := attr.Accept(v); err != nil {
            return nil, err
        }
    }
    a1 := v.cellsBuilder.Len()
    nodes := make([]*bufr.ValuedNode, a1-a0)
    for i := a0; i < a1; i++ {
        nodes[i-a0] = v.cellsBuilder.Cell(i).Node()
    }
    return nodes, nil
}
