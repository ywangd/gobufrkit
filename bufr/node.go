package bufr

import (
    "github.com/ywangd/gobufrkit/table"
    "fmt"
)

type Node interface {
    Members() []Node
    // Add a node as a child
    AddMember(node Node)
}

type Block []Node

func NewBlock() *Block {
    b := Block([]Node{})
    return &b
}

func (b *Block) Members() []Node {
    return *b
}

func (b *Block) AddMember(node Node) {
    *b = append([]Node(*b), node)
}

func (b *Block) Accept(visitor Visitor) error {
    return visitor.VisitBlock(b)
}

// ValuelessNode is a BUFR node that does not have an associated value.
//
// When it represents a delayed replication descriptor, the first member is
// the delayed replication factor. Client of the object can figure it out
// from the node's Descriptor.
type ValuelessNode struct {
    Descriptor table.Descriptor
    members    []Node
}

func (n *ValuelessNode) Members() []Node {
    return n.members
}

func (n *ValuelessNode) AddMember(node Node) {
    n.members = append(n.members, node)
}

func (n *ValuelessNode) Accept(visitor Visitor) error {
    return visitor.VisitValuelessNode(n)
}

// ValuedNode is a BUFR node that has an associated value.
// The value is represented by an Index to the value list
type ValuedNode struct {
    Descriptor  table.Descriptor
    Index       int
    PackingInfo *PackingInfo

    MinValue  interface{}
    NbitsDiff int

    members []Node // attributes
}

func (n *ValuedNode) Accept(visitor Visitor) error {
    return visitor.VisitValuedNode(n)
}

func (n *ValuedNode) Members() []Node {
    return n.members
}

func (n *ValuedNode) AddMember(node Node) {
    n.members = append(n.members, node)
}

func (n *ValuedNode) String() string {
    return fmt.Sprintf("[%d] %v (%d)", n.Index, n.Descriptor, n.PackingInfo.Nbits)
}