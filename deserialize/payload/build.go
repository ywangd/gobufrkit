package payload

import (
    "github.com/ywangd/gobufrkit/bufr"
    "github.com/ywangd/gobufrkit/deserialize/ast"
    "github.com/pkg/errors"
    "github.com/ywangd/gobufrkit/table"
)

// assocPair is a wrapper of a pair of data representing the number of bits
// and the associated field significance node.
type assocPair struct {
    Nbits int
    Node  bufr.Node
}

// assocPairs provides methods for managing a slice of assocPair
type assocPairs struct {
    pairs []*assocPair
}

func (p *assocPairs) Push(nbits int) {
    p.pairs = append(p.pairs, &assocPair{Nbits: nbits})
}

func (p *assocPairs) Pop() {
    p.pairs = p.pairs[:len(p.pairs)-1]
}

// SetNode is a no-op if there is no assocPair
func (p *assocPairs) SetNode(node bufr.Node) {
    if len(p.pairs) > 0 {
        p.pairs[len(p.pairs)-1].Node = node
    }
}

func (p *assocPairs) Pairs() []*assocPair {
    return p.pairs
}

// buildBlocks builds a block of replicated nodes
func buildBlock(v *DesVisitor, members []ast.Node) error {
    v.treeBuilder.Push(bufr.NewBlock())
    defer v.treeBuilder.Pop()
    for _, m := range members { // loop of each replicated block
        if err := m.Accept(v); err != nil {
            return errors.Wrap(err, "cannot process replicated block")
        }
    }
    return nil
}

// buildAssocNodes builds a slice of nodes that represents the list of associated field node
func buildAssocNodes(v *DesVisitor, descriptor table.Descriptor) ([]bufr.Node, error) {
    if descriptor.X() == 31 || len(v.assocPairs.Pairs()) == 0 {
        return nil, nil
    }
    nodes := []bufr.Node{}
    for _, p := range v.assocPairs.Pairs() {
        info := &bufr.PackingInfo{Unit: table.NONNEG_CODE, Nbits: p.Nbits}
        node, err := buildValuedNodeWithInfo(v, &table.DecorateDescriptor{
            Descriptor: descriptor, Initial: 'A', Name: "ASSOCIATED FIELD"}, info)
        if err != nil {
            return nil, errors.Wrap(err, "cannot build associated field node")
        }
        node.AddMember(p.Node)
        nodes = append(nodes, node)
    }
    return nodes, nil
}

// buildValuedNode builds a ValuedNode for the given descriptor.
// The packing info is calculated for the given descriptor.
// It calls the helper buildValueNodeWithInfo to do the actual building work.
func buildValuedNode(v *DesVisitor, descriptor table.Descriptor) (*bufr.ValuedNode, error) {
    info, err := calcPackingInfo(v, descriptor)
    if err != nil {
        return nil, errors.Wrap(err, "cannot calculate packing info")
    }
    return buildValuedNodeWithInfo(v, descriptor, info)
}

// buildValueNodeWithInfo unpack value(s) of the given descriptor, assemble the ValuedNode
// and also call treeBuilder and cellsBuilder to add the node and value(s).
func buildValuedNodeWithInfo(v *DesVisitor, descriptor table.Descriptor, info *bufr.PackingInfo) (*bufr.ValuedNode, error) {
    val, err := v.unpacker.Unpack(info)
    if err != nil {
        return nil, errors.Wrap(err, "cannot unpack value")
    }
    node := &bufr.ValuedNode{Descriptor: descriptor, PackingInfo: info}
    v.treeBuilder.Add(node)
    v.cellsBuilder.Add(node, val)
    return node, nil
}

// buildZeroNode adds a zero valued node if in compatible mode
func buildZeroNode(v *DesVisitor, descriptor table.Descriptor) error {
    if v.config.Compatible {
        info := &bufr.PackingInfo{Unit: table.NONNEG_CODE, Nbits: 0}
        _, err := buildValuedNodeWithInfo(v, descriptor, info)
        return err
    }
    v.treeBuilder.Add(&bufr.ValuelessNode{Descriptor: descriptor})
    return nil
}
