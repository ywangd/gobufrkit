package table

import (
    "testing"
    assert2 "github.com/seanpont/assert"
    "fmt"
)

func getTableGroup() TableGroup {
    b := &B{entries: map[ID]*Bentry{
        1001:  {name: "001001"},
        1002:  {name: "001002"},
        4003:  {name: "004003"},
        4004:  {name: "004004"},
        4005:  {name: "004005"},
        31001: {name: "031001"},
    }}

    d := &D{entries: map[ID]*Dentry{
        301001: {"301001", []ID{1001, 1002}},
        301013: {"301013", []ID{4003, 4004, 4005}},
    }}

    return &SingleTableGroup{b: b, d: d}
}

func TestUnexpandedTemplate_Expand(t *testing.T) {
    assert := assert2.Assert(t)
    group := getTableGroup()

    ut := UnexpandedTemplate{ids: []ID{301001, 301013}}
    et, _ := ut.Expand(group)
    assert.Equal(et.Dump(), `301001 301001
    001001 001001
    001002 001002
301013 301013
    004003 004003
    004004 004004
    004005 004005`)
}

func TestUnexpandedTemplate_Expand2(t *testing.T) {
    assert := assert2.Assert(t)
    group := getTableGroup()

    ut := UnexpandedTemplate{ids: []ID{103002, 1001, 1002, 4003}}
    et, _ := ut.Expand(group)
    assert.Equal(et.Dump(), `103002 103002
    001001 001001
    001002 001002
    004003 004003`)
}

func TestUnexpandedTemplate_Expand3(t *testing.T) {
    assert := assert2.Assert(t)
    group := getTableGroup()

    ut := UnexpandedTemplate{ids: []ID{103002, 4004, 301001, 4003}}
    et, _ := ut.Expand(group)

    assert.Equal(et.Dump(), `103002 103002
    004004 004004
    301001 301001
        001001 001001
        001002 001002
    004003 004003`)
}

func TestUnexpandedTemplate_Expand4(t *testing.T) {
    assert := assert2.Assert(t)
    group := getTableGroup()

    ut := UnexpandedTemplate{ids: []ID{103000, 31001, 301001, 301013, 1001}}
    et, _ := ut.Expand(group)
    assert.Equal(et.Dump(), `103000 103000
....031001 031001
    301001 301001
        001001 001001
        001002 001002
    301013 301013
        004003 004003
        004004 004004
        004005 004005
    001001 001001`)
}

func TestExpandedTemplate_Walk(t *testing.T) {
    assert := assert2.Assert(t)
    group := getTableGroup()

    ut := UnexpandedTemplate{ids: []ID{103000, 31001, 301001, 301013, 1001}}
    et, _ := ut.Expand(group)
    ins := make(chan int, 1)
    ins <- 2

    events := []string{
        "WE_BEGIN_TEMPLATE,<nil>",
        "WE_REPLICATION_DESCRIPTOR,103000",
        "WE_FACTOR,031001",
        "WE_BEGIN_REPLICATION,<nil>",
        "WE_BEGIN_BLOCK,<nil>",
        "WE_SEQUENCE_DESCRIPTOR,301001",
        "WE_BEGIN_SEQUENCE,<nil>",
        "WE_ELEMENT_DESCRIPTOR,001001",
        "WE_ELEMENT_DESCRIPTOR,001002",
        "WE_END_SEQUENCE,<nil>",
        "WE_SEQUENCE_DESCRIPTOR,301013",
        "WE_BEGIN_SEQUENCE,<nil>",
        "WE_ELEMENT_DESCRIPTOR,004003",
        "WE_ELEMENT_DESCRIPTOR,004004",
        "WE_ELEMENT_DESCRIPTOR,004005",
        "WE_END_SEQUENCE,<nil>",
        "WE_ELEMENT_DESCRIPTOR,001001",
        "WE_END_BLOCK,<nil>",
        "WE_BEGIN_BLOCK,<nil>",
        "WE_SEQUENCE_DESCRIPTOR,301001",
        "WE_BEGIN_SEQUENCE,<nil>",
        "WE_ELEMENT_DESCRIPTOR,001001",
        "WE_ELEMENT_DESCRIPTOR,001002",
        "WE_END_SEQUENCE,<nil>",
        "WE_SEQUENCE_DESCRIPTOR,301013",
        "WE_BEGIN_SEQUENCE,<nil>",
        "WE_ELEMENT_DESCRIPTOR,004003",
        "WE_ELEMENT_DESCRIPTOR,004004",
        "WE_ELEMENT_DESCRIPTOR,004005",
        "WE_END_SEQUENCE,<nil>",
        "WE_ELEMENT_DESCRIPTOR,001001",
        "WE_END_BLOCK,<nil>",
        "WE_END_REPLICATION,<nil>",
        "WE_END_TEMPLATE,<nil>",
    }

    i := 0
    for e := range et.Walk(ins) {
        if e.Descriptor != nil {
            assert.Equal(fmt.Sprintf("%v,%v", e.Code, e.Descriptor.Id()), events[i])
        } else {
            assert.Equal(fmt.Sprintf("%v,<nil>", e.Code), events[i])
        }
        i += 1
    }
    close(ins)
}
