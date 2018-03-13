package parser

import (
    "testing"
    assert2 "github.com/seanpont/assert"
    "github.com/ywangd/gobufrkit/table"
)

func getKeeper() *idsKeeper {
    return newIdsKeeper([]table.ID{
        236000,
        101000,
        31002,
        31031,
        1031,
        1032,
        101000,
        31002,
        33007,
    })
}

func TestIdsKeeper_TakeWhile(t *testing.T) {
    assert := assert2.Assert(t)

    ids := getKeeper().takeWhile(func(id table.ID) bool {
        return id == table.ID(31031)
    })

    assert.Equal(len(ids), 4)
    assert.Equal(ids[0], table.ID(236000))
    assert.Equal(ids[3], table.ID(31031))
}

func TestIdsKeeper_TakeTill(t *testing.T) {
    assert := assert2.Assert(t)

    ids := getKeeper().takeTill(func(id table.ID) bool {
        return id == table.ID(1032)
    })

    assert.Equal(len(ids), 5)
    assert.Equal(ids[0], table.ID(236000))
    assert.Equal(ids[4], table.ID(1031))
}

func TestParser_State(t *testing.T) {
    assert := assert2.Assert(t)
    parser := NewParser(nil)

    assert.False(parser.getState(stateOpDataNotPresent), "state should be false")
    parser.setState(stateOpDataNotPresent)
    assert.True(parser.getState(stateOpDataNotPresent), "state should be true")
    parser.unsetState(stateOpDataNotPresent)
    assert.False(parser.getState(stateOpDataNotPresent), "state should be false")
}
