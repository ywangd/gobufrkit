//go:generate stringer -type=EventCode
package table

import (
    "fmt"
    "strings"
    "encoding/json"
)

const (
    INDENT = "    "
    INDOTS = "...."
)

// EventCode is the code of Event generated while walking an ExpandedTemplate.
type EventCode int

const (
    WE_ELEMENT_DESCRIPTOR     EventCode = iota
    WE_REPLICATION_DESCRIPTOR
    WE_OPERATOR_DESCRIPTOR
    WE_SEQUENCE_DESCRIPTOR
    WE_FACTOR
    WE_BEGIN_SEQUENCE
    WE_END_SEQUENCE
    WE_BEGIN_REPLICATION
    WE_END_REPLICATION
    WE_BEGIN_BLOCK
    WE_END_BLOCK
    WE_BEGIN_TEMPLATE
    WE_END_TEMPLATE
)

// Event is the event generated while walking an ExpandedTemplate.
// It always has a code and potentially an associated Descriptor.
// If the event has no associated Descriptor, the field is nil.
type Event struct {
    Code       EventCode
    Descriptor Descriptor
}

// UnexpandedTemplate is the basic representation of template that is essentially
// a list of descriptor IDs
type UnexpandedTemplate struct {
    ids   []ID
    fbits int
    xbits int
    ybits int
}

func NewUnexpandedTemplate(ids []ID, fbits, xbits, ybits int) *UnexpandedTemplate {
    return &UnexpandedTemplate{ids: ids, fbits: fbits, xbits: xbits, ybits: ybits}
}

func (ut *UnexpandedTemplate) String() string {
    return fmt.Sprintf("%v", ut.ids)
}

func (ut *UnexpandedTemplate) Ids() []ID {
    return ut.ids
}

func (ut *UnexpandedTemplate) Fbits() int {
    return ut.fbits
}

func (ut *UnexpandedTemplate) Xbits() int {
    return ut.xbits
}

func (ut *UnexpandedTemplate) Ybits() int {
    return ut.ybits
}

// Expand creates an ExpandedTemplate by expanding the UnexpandedTemplate
func (ut *UnexpandedTemplate) Expand(group TableGroup) (*ExpandedTemplate, error) {

    nodes, err := expand(ut.ids, group)
    if err != nil {
        return nil, err
    }
    return &ExpandedTemplate{nodes}, nil
}

func (ut *UnexpandedTemplate) MarshalJSON() ([]byte, error) {
    return json.Marshal(ut.ids)
}

// ExpandedTemplate is the expanded form of an UnexpandedTemplate. It is structure
// is hierarchical to reflect the presence of sequence and replication descriptors.
type ExpandedTemplate struct {
    members []*Node
}

// Walk walks through the ExpandedTemplate and generating Events along the way.
// The Events fully reflects the hierarchical structure of the ExpandedTemplate.
// This should help the consumers to rebuild necessary hierarchy by only relying
// on the Event.
// The argument ins is used for the walker to get information about replication
// blocks. If ins is nil, all replication blocks will only be walked through once
// regardless of its actual number of repeats. This is useful for functions like
// dumping the structure of UnexpandedTemplate. When ins is NOT nil, the number of
// repeats will be read from ins when a delayed replication block is encountered.
// For fixed replication block, the number of repeats is the X value of the
// replication descriptor as described by the Code Manual.
//
// The function is defined to be a method of ExpandedTemplate instead of a separate
// Walker type is because there is only one way to walk the ExpandedTemplate.
func (et *ExpandedTemplate) Walk(ins chan int) chan *Event {
    outs := make(chan *Event)
    go func() {
        outs <- &Event{WE_BEGIN_TEMPLATE, nil}
        for _, n := range et.members {
            for e := range n.Walk(ins) {
                outs <- e
            }
        }
        outs <- &Event{WE_END_TEMPLATE, nil}
        close(outs)
    }()
    return outs
}

// Dump generates an user-friendly presentation of the hierarchical structure
func (et *ExpandedTemplate) Dump() string {
    lines := make([]string, len(et.members))
    for i, n := range et.members {
        lines[i] = n.Dump()
    }
    return strings.Join(lines, "\n")
}

// Node is a container of a descriptor that helps to form the ExpandedTemplate.
// In addition to holding a reference to the descriptor, it also has potential
// references to a factor node (for delayed replication) and a list of member
// nodes (for replication and sequence descriptors)
type Node struct {
    descriptor Descriptor
    factor     *Node
    members    []*Node
}

// Walk walks the node and any of its factor and members and generating
// Event along the way. It is called by ExpandedTemplate.Walk
func (n *Node) Walk(ins chan int) chan *Event {
    outs := make(chan *Event)
    go func() {
        switch n.descriptor.F() {
        case F_ELEMENT:
            outs <- &Event{WE_ELEMENT_DESCRIPTOR, n.descriptor}
        case F_REPLICATION:
            outs <- &Event{WE_REPLICATION_DESCRIPTOR, n.descriptor}
            nrepeats := 1
            if n.factor != nil {
                outs <- &Event{WE_FACTOR, n.factor.descriptor}
                if ins != nil {
                    nrepeats = <-ins
                }
            } else if ins != nil {
                nrepeats = n.descriptor.Y()
            }
            walkReplication(ins, outs, n.members, nrepeats)

        case F_OPERATOR:
            outs <- &Event{WE_OPERATOR_DESCRIPTOR, n.descriptor}

        case F_SEQUENCE:
            outs <- &Event{WE_SEQUENCE_DESCRIPTOR, n.descriptor}
            walkSequence(ins, outs, n.members)

        }
        close(outs)
    }()
    return outs
}

// Dump generates an user-friendly representation of a Node showing the descriptor
// itself and any factor and members.
func (n *Node) Dump() string {
    lines := []string{}
    indents := []string{}
    for e := range n.Walk(nil) {
        switch e.Code {
        case WE_ELEMENT_DESCRIPTOR, WE_REPLICATION_DESCRIPTOR, WE_OPERATOR_DESCRIPTOR, WE_SEQUENCE_DESCRIPTOR:
            lines = append(lines, fmt.Sprintf("%s%s", strings.Join(indents, ""), e.Descriptor))
        case WE_FACTOR:
            lines = append(lines, fmt.Sprintf("%s%s%s", strings.Join(indents, ""), INDOTS, e.Descriptor))
        case WE_BEGIN_SEQUENCE, WE_BEGIN_BLOCK:
            indents = append(indents, INDENT)
        case WE_END_SEQUENCE, WE_END_BLOCK:
            indents = indents[0:len(indents)-1]
        }
    }
    return strings.Join(lines, "\n")
}

// The walkSequence is a helper function to walk through a list of Node
func walkSequence(ins chan int, outs chan *Event, members []*Node) {
    outs <- &Event{WE_BEGIN_SEQUENCE, nil}
    for _, n := range members {
        for e := range n.Walk(ins) {
            outs <- e
        }
    }
    outs <- &Event{WE_END_SEQUENCE, nil}
}

// The walkReplication is a helper function to walk through a list of Node
// a number of times.
func walkReplication(ins chan int, outs chan *Event, members []*Node, nrepeats int) {
    outs <- &Event{WE_BEGIN_REPLICATION, nil}
    for i := 0; i < nrepeats; i++ {
        outs <- &Event{WE_BEGIN_BLOCK, nil}
        for _, n := range members {
            for e := range n.Walk(ins) {
                outs <- e
            }
        }
        outs <- &Event{WE_END_BLOCK, nil}
    }
    outs <- &Event{WE_END_REPLICATION, nil}
}

// expand helps to expand an UnexpandedTemplate by recursively descend into
// sequence descriptors and group replicated blocks
//
// The function takes a slice of ID and expand its members using the given TableGroup.
func expand(ids []ID, group TableGroup) ([]*Node, error) {

    nodes := []*Node{}
    count := 0
    for count < len(ids) {
        d, err := group.Lookup(ids[count])
        if err != nil {
            return nil, err
        }
        count += 1
        n := &Node{descriptor: d}
        nodes = append(nodes, n)
        switch d.F() {
        case F_SEQUENCE:
            members, err := expand(d.Entry().(*Dentry).Members, group)
            if err != nil {
                return nil, err
            }
            n.members = members

        case F_REPLICATION:
            if d.Y() == 0 { // delayed replication factor
                d, err := group.Lookup(ids[count])
                if err != nil {
                    return nil, err
                }
                count += 1
                n.factor = &Node{descriptor: d}
            }
            replicated := d.X()
            members, err := expand(ids[count: count+replicated], group)
            if err != nil {
                return nil, err
            }
            count += replicated
            n.members = members
        }
    }
    return nodes, nil
}
