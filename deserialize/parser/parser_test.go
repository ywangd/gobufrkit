package parser_test

import (
    "testing"
    assert2 "github.com/seanpont/assert"
    "github.com/ywangd/gobufrkit/table"
    "github.com/ywangd/gobufrkit/deserialize/parser"
    "github.com/ywangd/gobufrkit/deserialize/ast"
    "os"
)

func TestParser_Parse(t *testing.T) {
    assert := assert2.Assert(t)

    tableGroup, err := table.NewSingleTableGroup(
        "../../_definitions/tables",
        0, 0, 0, 28)
    assert.Nil(err)

    p := parser.NewParser(tableGroup)

    ut := table.NewUnexpandedTemplate([]table.ID{
        301001, // sequence descriptor
        102001, // fixed replication
        4001,
        4002,
        103000, // delayed replication
        31001,
        5001,
        5002,
        6001,
        221003, // data not present
        5001,
        5002,
        6001,
        222000, // QA info follows
        236000, // define reusable bitmap
        101000,
        31002,
        31031,
        1031,
        1032,
        101000,
        31002,
        33007,
        237255, // cancel bitmap
        223000, // first order stats
        101000,  // ad-hoc bitmap
        31002,
        31031,
        1031,
        1032,
        8023,
        223255,
        223255,
        204003, // associated field
        31021,
        206050, // skip local descriptor
        5003,   // non-existing descriptor
        31021,  // change meaning of current associated field
        204000, // cancel associated field
        203012, // new refval of 12 bits
        1031,
        1032,
        203255,
    }, 0, 0, 0)
    tree, err := p.Parse(ut)
    assert.Nil(err)

    visitor := ast.DumpVisitor(os.Stdout)

    assert.Nil(tree.Accept(visitor))

}
