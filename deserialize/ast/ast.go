// Package ast provides building blocks for constructing the tree of parsed template.
// The ast is the parsed template.
package ast

import (
    "github.com/ywangd/gobufrkit/table"
)

// Node is the interface that all nodes of a parse template tree must implement.
type Node interface {
    Descriptor() table.Descriptor
    SetMembers([]Node)  // TODO: change to AddMembers with varargs?
    Members() []Node
    Accept(visitor Visitor) error
}

// BaseNode is a basic implementation of the Node interface.
type BaseNode struct {
    // descriptor associated to the node
    descriptor table.Descriptor
    // members nodes
    members []Node
}

func NewBaseNode(descriptor table.Descriptor) *BaseNode {
    return &BaseNode{descriptor: descriptor}
}

func (n *BaseNode) Descriptor() table.Descriptor {
    return n.descriptor
}

func (n *BaseNode) Members() []Node {
    return n.members
}

func (n *BaseNode) SetMembers(members []Node) {
    n.members = members
}

func (n *BaseNode) Accept(visitor Visitor) error {
    return visitor.VisitNode(n)
}

type ElementNode struct {
    *BaseNode
    NotPresent bool
}

func (n *ElementNode) Accept(visitor Visitor) error {
    return visitor.VisitElementNode(n)
}

// E031021Node represents a standalone associated field significance descriptor.
// Standalone means it is not immediately after a 204YYY associated field operator.
type E031021Node struct {
    *BaseNode
}

func (n *E031021Node) Accept(visitor Visitor) error {
    return visitor.VisitE031021Node(n)
}

type SequenceNode struct {
    *BaseNode
}

func (n *SequenceNode) Accept(visitor Visitor) error {
    return visitor.VisitSequenceNode(n)
}

type FixedReplicationNode struct {
    *BaseNode
}

func (n *FixedReplicationNode) Accept(visitor Visitor) error {
    return visitor.VisitFixedReplicationNode(n)
}

// The first member is the delayed replication factor
type DelayedReplicationNode struct {
    *BaseNode
}

func (n *DelayedReplicationNode) Accept(visitor Visitor) error {
    return visitor.VisitDelayedReplicationNode(n)
}

// OpNbitsOffsetNode represents operator 201YYY
type OpNbitsOffsetNode struct {
    *BaseNode
}

func (n *OpNbitsOffsetNode) Accept(visitor Visitor) error {
    return visitor.VisitOpNbitsOffsetNode(n)
}

// OpScaleOffsetNode represents operator 202YYY
type OpScaleOffsetNode struct {
    *BaseNode
}

func (n *OpScaleOffsetNode) Accept(visitor Visitor) error {
    return visitor.VisitOpScaleOffsetNode(n)
}

// OpNewRefvalNode represents operator 203YYY
type OpNewRefvalNode struct {
    *BaseNode
    // Its only member is the descriptor which the new refval will be applied to
}

func (n *OpNewRefvalNode) Accept(visitor Visitor) error {
    return visitor.VisitOpNewRefvalNode(n)
}

// OpAssocFieldNode represents operator 204YYY
type OpAssocFieldNode struct {
    *BaseNode
    // 204YYY (except 204000) has an only member of associated field significance
}

func (n *OpAssocFieldNode) Accept(visitor Visitor) error {
    return visitor.VisitOpAssocFieldNode(n)
}

// OpInsertStringNode represents operator 205YYY
type OpInsertStringNode struct {
    *BaseNode
}

func (n *OpInsertStringNode) Accept(visitor Visitor) error {
    return visitor.VisitOpInsertStringNode(n)
}

// OpInsertStringNode represents operator 206YYY
type OpSkipLocalNode struct {
    *BaseNode
    // its only member is the local descriptor to be skipped
}

func (n *OpSkipLocalNode) Accept(visitor Visitor) error {
    return visitor.VisitOpSkipLocalNode(n)
}

// 207YYY
type OpModifyPackingNode struct {
    *BaseNode
}

func (n *OpModifyPackingNode) Accept(visitor Visitor) error {
    return visitor.VisitOpModifyPackingNode(n)
}

// 208YYY
type OpSetStringLengthNode struct {
    *BaseNode
}

func (n *OpSetStringLengthNode) Accept(visitor Visitor) error {
    return visitor.VisitOpSetStringLengthNode(n)
}

// OpDataNotPresentNode represents operator 221YYY
type OpDataNotPresentNode struct {
    *BaseNode
}

func (n *OpDataNotPresentNode) Accept(visitor Visitor) error {
    return visitor.VisitOpDataNotPresentNode(n)
}

// OpAssessmentNode contains a sequence of descriptors related a info or stats session.
// This includes QA Info 222000, Stats 224000, 225000, substitution
// 222000, replacement 232000
type OpAssessmentNode struct {
    *BaseNode
    Bitmap Node
    // These are the descriptors sandwiched between bitmap and class 33 QA descriptors.
    Attrs []Node
}

func (n *OpAssessmentNode) Accept(visitor Visitor) error {
    return visitor.VisitOpAssessmentNode(n)
}

// OpMarkerNode represents all the marker descriptors, e.g. 223255, 224255
type OpMarkerNode struct {
    *BaseNode
}

func (n *OpMarkerNode) Accept(visitor Visitor) error {
    return visitor.VisitOpMarkerNode(n)
}

type OpCancelBackRefNode struct {
    *BaseNode
}

func (n *OpCancelBackRefNode) Accept(visitor Visitor) error {
    return visitor.VisitOpCancelBackRefNode(n)
}

// BitmapNode represents either a definition or a recall of a bitmap.
type BitmapNode struct {
    *BaseNode
    // Descriptor is first descriptor that begins the bitmap.
    // It can be either 236000, 237000 or nil if it is an ad-hoc bitmap definition,
    // i.e. not for reuse
    // Often a bitmap is defined within other section such as quality info, i.e.
    // 222000, 236000, ...
    //
    // However, it is also possible to have a standalone bitmap definition
    // that begins with a 236000, i.e. 236000, ...
}

func (n *BitmapNode) Accept(visitor Visitor) error {
    return visitor.VisitBitmapNode(n)
}

type OpCancelBitmapNode struct {
    *BaseNode
}

func (n *OpCancelBitmapNode) Accept(visitor Visitor) error{
    return visitor.VisitOpCancelBitmapNode(n)
}