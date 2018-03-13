package api

import (
    "github.com/Shopify/go-lua"
    "github.com/ywangd/gobufrkit/bufr"
)

const SECTION_META_TABLE = "gobufrkit.section"

type ObjSection struct {
}

func (obj *ObjSection) getField(state *lua.State) int {
    section := state.ToUserData(1).(*bufr.Section)
    name, _ := state.ToString(2)
    pushField(state, section.FieldByName(name))
    return 1
}

func (obj *ObjSection) getNumberOfFields(state *lua.State) int {
    section := state.ToUserData(1).(*bufr.Section)
    state.PushInteger(len(section.Fields()))
    return 1
}

func (obj *ObjSection) registerGribSectionType(state *lua.State) {
    lua.NewMetaTable(state, SECTION_META_TABLE)

    // metatable.__index = metatable
    state.PushValue(-1)
    state.SetField(-2, "__index")

    lua.SetFunctions(state, []lua.RegistryFunction{
        {Name: "getField", Function: obj.getField},
        {Name: "getNumberOfFields", Function: obj.getNumberOfFields},
    }, 0)
}
