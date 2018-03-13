package api

import (
    "github.com/Shopify/go-lua"
    "fmt"
    "github.com/ywangd/gobufrkit/bufr"
)

const FIELD_META_TABLE = "gobufrkit.field"

type ObjField struct {
}

// Provide nice to read string representation of the field object for lua
// TODO: Improve the string format
func (obj *ObjField) fieldToString(state *lua.State) int {
    field := state.ToUserData(1).(*bufr.Field)
    state.PushString(fmt.Sprintf("%+v", field))
    return 1
}

func (obj *ObjField) value(state *lua.State) int {
    field := state.ToUserData(1).(*bufr.Field)
    err := pushFieldValue(state, field)
    if err != nil {
        state.PushString(err.Error())
        state.Error()
        return 0
    }
    return 1
}

func (obj *ObjField) isMissing(state *lua.State) int {
    field := state.ToUserData(1).(*bufr.Field)
    state.PushBoolean(field.IsMissing())
    return 1
}

func (obj *ObjField) registerSectionFieldType(state *lua.State) {
    lua.NewMetaTable(state, FIELD_META_TABLE)

    // metatable.__index = metatable
    state.PushValue(-1)
    state.SetField(-2, "__index")

    lua.SetFunctions(state, []lua.RegistryFunction{
        {Name: "__tostring", Function: obj.fieldToString},
        {Name: "value", Function: obj.value},
        {Name: "isMissing", Function: obj.isMissing},
    }, 0)
}
