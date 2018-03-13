package payload

import (
    "github.com/ywangd/gobufrkit/bufr"
    "fmt"
    "github.com/ywangd/gobufrkit/debug"
)

// TreeBuilder provides operations for building the hierarchical BUFR nodes.
type TreeBuilder struct {
    Verbose bool

    node  bufr.Node
    stack []bufr.Node
}

func NewTreeBuilder(node bufr.Node) *TreeBuilder {
    return &TreeBuilder{node: node}
}

// Add given node as a child to the current parent node
func (tb *TreeBuilder) Add(node bufr.Node) {
    if debug.DEBUG {
        fmt.Println("++++++++++++++ START ADD NODE ++++++++++++++")
        defer fmt.Println("++++++++++++++ DONE ADD NODE ++++++++++++++")
    }
    if debug.DEBUG {
        fmt.Println(node)
    }
    tb.node.AddMember(node)
}

// Add given node as a child to the current parent node and push it as the new parent
func (tb *TreeBuilder) Push(node bufr.Node) {
    if debug.DEBUG {
        fmt.Println(">>>>>>>>>>>>>> START PUSH NODE >>>>>>>>>>>>>>")
        defer fmt.Println(">>>>>>>>>>>>>> DONE PUSH NODE >>>>>>>>>>>>>>")
    }
    if debug.DEBUG {
        fmt.Println(len(tb.stack), node)
    }
    tb.stack = append(tb.stack, tb.node)
    tb.node = node
    if debug.DEBUG {
        fmt.Println(len(tb.stack))
    }
}

// Pop replace the current parent with its own parent
func (tb *TreeBuilder) Pop() {
    if debug.DEBUG {
        fmt.Println("<<<<<<<<<<<<<<< START POP NODE <<<<<<<<<<<<<<<")
        defer fmt.Println(" <<<<<<<<<<<<<<< DONE POP NODE <<<<<<<<<<<<<<<")
    }
    if debug.DEBUG {
        fmt.Println(len(tb.stack), tb.node)
    }
    tb.stack, tb.node = tb.stack[:len(tb.stack)-1], tb.stack[len(tb.stack)-1]
    if debug.DEBUG {
        fmt.Println(len(tb.stack))
    }
}

func (tb *TreeBuilder) Root() (bufr.Node, error) {
    if len(tb.stack) != 0 {
        return nil, fmt.Errorf("tree builder stack is not empty when root node is required")
    }
    return tb.node, nil
}
