// Package parser parses UnexpandedTemplate and construct a tree of the parsed template.
// The purpose of parsing is to reduce complexity of the deserializer.
// It does its job by:
//  * Localise contexts to remove the need of global state, e.g. skip local descriptor
//  * Reduce the need of switch/if/else by creating node of specific types
//  * Manage global state in finer and simpler control, e.g. bitmap definition
//
// The parsed template makes it easier to construct the hierarchical structure of the
// BUFR nodes. This is achieved by:
//  * similarities between the structures of parsed templates and BUFR nodes
//  * maintain the context for related descriptors, e.g. bitmap
//
// TODO: What is the best way to flag the same descriptor for different situation
//       e.g. 031021 standalone and non-standalone (following a 204YYY)
//       e.g. element descriptors inside or outside of a data not present session
package parser

import (
    "github.com/ywangd/gobufrkit/table"
    "github.com/ywangd/gobufrkit/deserialize/ast"
    "fmt"
    "math/bits"
)

// Parser creates a tree of ast.Node from an UnexpandedTemplate.
// The tree is created according to the structural information of the template.
// The goal of having this tree is to ease the deserialization process by
// providing context to every descriptor, which makes it easier to construct
// the hierarchical BUFR object.
type Parser struct {
    // The group of tables used for descriptor lookup
    tableGroup table.TableGroup
    // bit flags for states
    states uint
}

func NewParser(tableGroup table.TableGroup) *Parser {
    return &Parser{
        tableGroup: tableGroup,
    }
}

// Parse returns the root node of the parsed template tree.
func (p *Parser) Parse(ut *table.UnexpandedTemplate) (ast.Node, error) {
    return populateMembers(p, newIdsKeeper(ut.Ids()), ast.NewBaseNode(table.RootDescriptor))
}

func (p *Parser) setState(state uint) {
    p.states = p.states | state
}

func (p *Parser) unsetState(state uint) {
    p.states = p.states & ^state
}

func (p *Parser) getState(state uint) bool {
    return bits.OnesCount(p.states&state) == 1
}

// It is possible to add info about the associated fields here, i.e. number of
// associated fields and corresponding bits required. However, the significance
// values have to be deserialized from the data section. This means the information
// about associated fields still have to be tracked during deserialization process.
// Therefore it is not worthwhile to track it here.
// Also it is not right to store the significance value to the OpAssocFieldNode.
//  1. The structural tree should NOT be polluted with actual values.
//  2. If the OpAssocFieldNode is inside a replication, each of the replication
//     may have its own value which is unwieldy for a single node to store all
//     the values.
func (p *Parser) parseElementNode(keeper *idsKeeper) (ast.Node, error) {
    var (
        descriptor table.Descriptor
        err        error
    )
    if p.getState(stateOpSkipLocal) {
        // Create an ad-hoc local descriptor so it does not error out
        descriptor = table.NewLocalDescriptor(keeper.take())
    } else {
        if descriptor, err = p.tableGroup.Lookup(keeper.take()); err != nil {
            return nil, err
        }
    }
    switch descriptor.Id() {
    case table.ID_031021:
        return &ast.E031021Node{BaseNode: ast.NewBaseNode(descriptor)}, nil
    default:
        return &ast.ElementNode{BaseNode: ast.NewBaseNode(descriptor),
            NotPresent: notPresent(p, descriptor)}, nil
    }
}

func (p *Parser) parseSequenceNode(keeper *idsKeeper) (ast.Node, error) {
    descriptor, err := p.tableGroup.Lookup(keeper.take())
    if err != nil {
        return nil, err
    }
    return populateMembers(p, newIdsKeeper(descriptor.Entry().(*table.Dentry).Members),
        &ast.SequenceNode{BaseNode: ast.NewBaseNode(descriptor)})
}

func (p *Parser) parseDelayedReplicationNode(keeper *idsKeeper) (ast.Node, error) {
    descriptor, err := p.tableGroup.Lookup(keeper.take())
    if err != nil {
        return nil, err
    }
    return populateMembers(p, newIdsKeeper(keeper.takeN(descriptor.X()+1)),
        &ast.DelayedReplicationNode{BaseNode: ast.NewBaseNode(descriptor)})
}

func (p *Parser) parseFixedReplicationNode(keeper *idsKeeper) (ast.Node, error) {
    descriptor, err := p.tableGroup.Lookup(keeper.take())
    if err != nil {
        return nil, err
    }
    return populateMembers(p, newIdsKeeper(keeper.takeN(descriptor.X())),
        &ast.FixedReplicationNode{BaseNode: ast.NewBaseNode(descriptor)})
}

func (p *Parser) parseOperatorNode(keeper *idsKeeper) (ast.Node, error) {
    descriptor, err := p.tableGroup.Lookup(keeper.take())
    if err != nil {
        return nil, err
    }
    switch descriptor.Operator() {
    case table.OP_NBITS_OFFSET:
        return &ast.OpNbitsOffsetNode{BaseNode: ast.NewBaseNode(descriptor)}, nil

    case table.OP_SCALE_OFFSET:
        return &ast.OpScaleOffsetNode{BaseNode: ast.NewBaseNode(descriptor)}, nil

    case table.OP_NEW_REFVAL:
        return assembleOpNewRefvalNode(p, keeper, descriptor)

    case table.OP_ASSOCIATE_FIELD:
        return assembleOpAssocFieldNode(p, keeper, descriptor)

    case table.OP_INSERT_STRING:
        return &ast.OpInsertStringNode{BaseNode: ast.NewBaseNode(descriptor)}, nil

    case table.OP_SKIP_LOCAL:
        return assembleOpSkipLocalNode(p, keeper, descriptor)

    case table.OP_MODIFY_PACKING:
        return &ast.OpModifyPackingNode{BaseNode: ast.NewBaseNode(descriptor)}, nil

    case table.OP_SET_STRING_LENGTH:
        return &ast.OpSetStringLengthNode{BaseNode: ast.NewBaseNode(descriptor)}, nil

    case table.OP_DATA_NOT_PRESENT:
        return assembleOpDataNotPresentNode(p, keeper, descriptor)

    case table.OP_QUALITY_INFO:
        return assembleOpQaInfoNode(p, keeper, descriptor)

    case table.OP_SUBSTITUTION:
        return assembleOpSubstitutionNode(p, keeper, descriptor)

    case table.OP_FIRST_ORDER_STATS:
        return assembleOpFirstOrderStatsNode(p, keeper, descriptor)

    case table.OP_DIFFERENCE_STATS:
        return assembleOpDiffStatsNode(p, keeper, descriptor)

    case table.OP_REPLACEMENT:
        return assembleOpReplacementNode(p, keeper, descriptor)

    case table.OP_CANCEL_BACK_REFERENCE:
        return &ast.OpCancelBackRefNode{BaseNode: ast.NewBaseNode(descriptor)}, nil

    case table.OP_DEFINE_BITMAP:
        keeper.back()
        return p.parseBitmapNode(keeper)

    case table.OP_RECALL_BITMAP:
        if descriptor.Operand() == 255 {
            return &ast.OpCancelBitmapNode{BaseNode: ast.NewBaseNode(descriptor)}, nil
        }
        keeper.back()
        return p.parseBitmapNode(keeper)

    case table.OP_DEFINE_EVENT, table.OP_DEFINE_CONDITIONING_EVENT, table.OP_CATEGORICAL_VALUES:
        return nil, fmt.Errorf("operator not implemented: %v", descriptor)

    default:
        return nil, fmt.Errorf("unrecognised operator: %v", descriptor)
    }
}

// parseBitmapNode creates a BitmapNode by reading from the given keeper.
func (p *Parser) parseBitmapNode(keeper *idsKeeper) (ast.Node, error) {
    // The opening descriptor
    descriptor, err := p.tableGroup.Lookup(keeper.take())
    if err != nil {
        return nil, err
    }

    switch {
    case descriptor.Id() == table.ID_237000:
        // A bitmap recall requires no further processing
        return &ast.BitmapNode{BaseNode: ast.NewBaseNode(descriptor)}, nil

    case descriptor.F() == 1 || descriptor.Id() == table.ID_031031:
        // An ad-hoc bitmap definition, i.e. a definition follows directly after
        // descriptors such as QA Info (222000), first order stats (224000)
        keeper.back()
        descriptor = nil

    case descriptor.Id() == table.ID_236000:
        // A reusable bitmap definition

    default:
        return nil, fmt.Errorf("invalid bitmap definition: %v", descriptor)
    }

    // Get all following IDs that relate to the bitmap definition
    ids := keeper.takeWhile(func(id table.ID) bool {
        return id == table.ID_031031
    })
    return populateMembers(p, newIdsKeeper(ids), &ast.BitmapNode{BaseNode: ast.NewBaseNode(descriptor)})
}
