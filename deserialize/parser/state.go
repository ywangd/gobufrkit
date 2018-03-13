package parser

import (
    "github.com/ywangd/gobufrkit/table"
    "github.com/ywangd/gobufrkit/deserialize/ast"
    "fmt"
)

const (
    stateOpDataNotPresent = 1 << iota
    stateOpSkipLocal
)

// idsKeeper is a wrapper around a slice of IDs that provides tracking features.
type idsKeeper struct {
    // A slice of ID to work with
    ids []table.ID
    // index keeps track of the take read position to the ID slice
    index int
}

func newIdsKeeper(ids []table.ID) *idsKeeper {
    return &idsKeeper{ids: ids}
}

// eof tests whether the keeper has reached to the end of the slice
func (k *idsKeeper) eof() bool {
    return k.index >= len(k.ids)
}

// take returns the take ID and advance the reader head, i.e. index
// TODO: Error handling when there is no more ID to take?
func (k *idsKeeper) take() table.ID {
    id := k.ids[k.index]
    k.index += 1
    return id
}

// back decreases the reader head by one
func (k *idsKeeper) back() {
    k.index -= 1
}

// peek returns the next ID without advance the reader head, i.e. index
func (k *idsKeeper) peek() table.ID {
    id := k.take()
    k.back()
    return id
}

// takeN returns a slice of ID containing the next N IDs from the keeper
// TODO: Error handling when there are not enough IDs? This can be checked with LINT
// TODO: optimise by record start and stop index and make the slice in one go?
func (k *idsKeeper) takeN(n int) []table.ID {
    ids := make([]table.ID, n)
    for i := 0; i < n && !k.eof(); i++ {
        ids[i] = k.take()
    }
    return ids
}

// takeWhile takes next ID repeatedly while the given predicate returns true.
// An usual feature of this function is that it allows unmatched IDs to
// appear before the first match. That is the while match does not start
// till the first match appears.
// For an example, given a list of IDs like the follows:
// 236000, 101000, 31002, 31031, 1031, 1032
// A predicate to match 31031 will accept 236000, 101000 and 31002.
// The while match effectively does not start till the first 31031.
// Once the first 31031 is matched, it then works like a normal takeWhile,
// i.e. it will stop at the following 1031.
func (k *idsKeeper) takeWhile(predicate func(table.ID) bool) []table.ID {
    started := false
    ids := []table.ID{}
    for !k.eof() {
        id := k.take()
        match := predicate(id)
        if match && !started {
            started = true
        } else if !match && started {
            k.back()
            break
        }
        ids = append(ids, id)
    }
    return ids
}

// takeTill takes next ID repeated until the predicate returns true.
func (k *idsKeeper) takeTill(predicate func(table.ID) bool) []table.ID {
    ids := []table.ID{}
    for !k.eof() {
        id := k.take()
        if predicate(id) {
            k.back()
            break
        }
        ids = append(ids, id)
    }
    return ids
}

// processThrough process each ID available from the given keeper and call
// corresponding parse function.
func processThrough(p *Parser, keeper *idsKeeper) (nodes []ast.Node, err error) {
    var (
        node ast.Node
    )
    for !keeper.eof() {
        id := keeper.peek()

        switch {

        case id.F() == table.F_SEQUENCE:
            if node, err = p.parseSequenceNode(keeper); err != nil {
                return nil, err
            }

        case id.F() == table.F_REPLICATION && id.Y() == 0:
            if node, err = p.parseDelayedReplicationNode(keeper); err != nil {
                return nil, err
            }

        case id.F() == table.F_REPLICATION:
            if node, err = p.parseFixedReplicationNode(keeper); err != nil {
                return nil, err
            }

        case id.F() == table.F_OPERATOR:
            if node, err = p.parseOperatorNode(keeper); err != nil {
                return nil, err
            }

        case id.F() == table.F_ELEMENT:
            if node, err = p.parseElementNode(keeper); err != nil {
                return nil, err
            }

        default:
            return nil, fmt.Errorf("unknown ID: %v", id)
        }
        nodes = append(nodes, node)
    }
    return nodes, nil
}

// populateMembers sets members of the given node by processing through all IDs provided by the keeper.
func populateMembers(p *Parser, keeper *idsKeeper, node ast.Node) (ast.Node, error) {
    members, err := processThrough(p, keeper)
    if err != nil {
        return nil, err
    }
    node.SetMembers(members)
    return node, nil
}

func assembleOpNewRefvalNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    // TODO: Ensure new refval values from all subsets are equal for data of compressed packing
    //       Technically, they could be different, but this would create much trouble than its worth.
    //       In practice, its unlikely anybody would do that.
    if descriptor.Operand() == 255 {
        return &ast.OpNewRefvalNode{BaseNode: ast.NewBaseNode(descriptor)}, nil
    }
    return populateMembers(p, newIdsKeeper(keeper.takeTill(func(id table.ID) bool {
        return id == table.ID_203255
    })), &ast.OpNewRefvalNode{BaseNode: ast.NewBaseNode(descriptor)})
}

// Associated Field operator 204YYY must be immediately followed by a significance descriptor of 031021.
// The nodes affected by OpAssocFieldNode are not nested under due to a few reasons:
//  * Unlike replication, it is not very straightforward to get the number of affected nodes.
//  * Even if the nodes are nested under it, it does not help to localise the context.
//    That is, during deserialization, it is still necessary to check a global status
//    to decide whether associated fields present.
//  * For the BUFR object, its structure does NOT require affected nodes to be nested
//    under the operator node (unlike replication).
func assembleOpAssocFieldNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    node := &ast.OpAssocFieldNode{BaseNode: ast.NewBaseNode(descriptor)}
    if descriptor.Y() != 0 {
        return populateMembers(p, newIdsKeeper(keeper.takeN(1)), node)
    }
    return node, nil
}

func assembleOpSkipLocalNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    p.setState(stateOpSkipLocal)
    defer p.unsetState(stateOpSkipLocal)
    return populateMembers(p, newIdsKeeper(keeper.takeN(1)),
        &ast.OpSkipLocalNode{BaseNode: ast.NewBaseNode(descriptor)})
}

func assembleOpDataNotPresentNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    p.setState(stateOpDataNotPresent)
    defer p.unsetState(stateOpDataNotPresent)
    return populateMembers(p, newIdsKeeper(keeper.takeN(descriptor.Y())),
        &ast.OpDataNotPresentNode{BaseNode: ast.NewBaseNode(descriptor)})
}

// assembleOpQaInfoNode creates a StatsNode for QA Info session by using the given descriptor as
// the beginning descriptor and read IDs from the keeper when necessary
func assembleOpQaInfoNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    if descriptor.Id() != table.ID_222000 {
        return nil, fmt.Errorf("invalid operand: %v", descriptor)
    }
    return assembleAssessmentNode(p, keeper, descriptor,
        func(id table.ID) bool {
            return id.F() == 1 || (id.F() == 0 && id.X() == 33)
        },
        func(id table.ID) bool {
            return id.F() == 0 && id.X() == 33
        })
}

func assembleOpSubstitutionNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    switch descriptor.Id() {
    case table.ID_223000:
        return assembleAssessmentNode(p, keeper, descriptor,
            func(id table.ID) bool {
                return id.F() == 1 || id == table.ID_223255
            },
            func(id table.ID) bool {
                return id == table.ID_223255
            })
    case table.ID_223255:
        return &ast.OpMarkerNode{BaseNode: ast.NewBaseNode(descriptor)}, nil
    default:
        return nil, fmt.Errorf("invalid operand: %v", descriptor)
    }
}

func assembleOpFirstOrderStatsNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    switch descriptor.Id() {
    case table.ID_224000:
        return assembleAssessmentNode(p, keeper, descriptor,
            func(id table.ID) bool {
                return id.F() == 1 || id == table.ID_224255
            },
            func(id table.ID) bool {
                return id == table.ID_224255
            })
    case table.ID_224255:
        return &ast.OpMarkerNode{BaseNode: ast.NewBaseNode(descriptor)}, nil
    default:
        return nil, fmt.Errorf("invalid operand: %v", descriptor)
    }
}

func assembleOpDiffStatsNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    switch descriptor.Id() {
    case table.ID_225000:
        return assembleAssessmentNode(p, keeper, descriptor,
            func(id table.ID) bool {
                return id.F() == 1 || id == table.ID_225255
            },
            func(id table.ID) bool {
                return id == table.ID_225255
            })
    case table.ID_225255:
        return &ast.OpMarkerNode{BaseNode: ast.NewBaseNode(descriptor)}, nil
    default:
        return nil, fmt.Errorf("invalid operand: %v", descriptor)
    }
}

func assembleOpReplacementNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor) (ast.Node, error) {
    switch descriptor.Id() {
    case table.ID_232000:
        return assembleAssessmentNode(p, keeper, descriptor,
            func(id table.ID) bool {
                return id.F() == 1 || id == table.ID_232255
            },
            func(id table.ID) bool {
                return id == table.ID_232255
            })
    case table.ID_232255:
        return &ast.OpMarkerNode{BaseNode: ast.NewBaseNode(descriptor)}, nil
    default:
        return nil, fmt.Errorf("invalid operand: %v", descriptor)
    }
}

// Helper function to assemble a OpAssessmentNode of different descriptor
// attrPredicate is a takeTill predicate for finding all sandwiched attr descriptors
// membersPredicate is a takeWhile predicate for finding all member descriptors
func assembleAssessmentNode(p *Parser, keeper *idsKeeper, descriptor table.Descriptor,
    attrPredicate, membersPredicate func(table.ID) bool) (ast.Node, error) {
    // Create the associated bitmap
    bitmapNode, err := p.parseBitmapNode(keeper)
    if err != nil {
        return nil, err
    }

    // Get all sandwiched descriptors between bitmap and actual data
    attrs, err := processAssessmentNodeAttrs(p, keeper, attrPredicate)
    if err != nil {
        return nil, err
    }

    // Get all the actual class 33 QA descriptors and use them as node members
    return processAssessmentNodeMembers(p, keeper, membersPredicate,
        &ast.OpAssessmentNode{
            BaseNode: ast.NewBaseNode(descriptor),
            Bitmap:   bitmapNode,
            Attrs:    attrs,
        })
}

func processAssessmentNodeAttrs(p *Parser, keeper *idsKeeper,
    predicate func(table.ID) bool) ([]ast.Node, error) {
    // There are often extra descriptors sandwiched between the bitmap section and QA descriptors,
    // e.g. 222000, 236000, 101000, 031002, 031031, 001031, 001032, 101000, 031002, 033007
    // Note the descriptors 001031 and 001032 are the extra descriptors.
    ids := keeper.takeTill(predicate)
    return processThrough(p, newIdsKeeper(ids))
}

//
func processAssessmentNodeMembers(p *Parser, keeper *idsKeeper,
    predicate func(table.ID) bool, node ast.Node) (ast.Node, error) {
    // Get all the actual class 33 QA descriptors or marker nodes
    ids := keeper.takeWhile(predicate)
    return populateMembers(p, newIdsKeeper(ids), node)
}

// Test whether an element descriptor should be not present due to 221YYY
func notPresent(p *Parser, descriptor table.Descriptor) bool {
    if !p.getState(stateOpDataNotPresent) {
        return false
    }
    x := descriptor.X()
    return x < 1 && x > 9 && x != 31
}
