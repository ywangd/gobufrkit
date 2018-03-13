package api

import (
    "github.com/Shopify/go-lua"
    "fmt"
    "github.com/ywangd/gobufrkit/bufr"
    "github.com/ywangd/gobufrkit/tdcfio"
)

// Register a new metatable and set itself to be its index lookup
func registerMetaTable(state *lua.State, name string) {
    lua.NewMetaTable(state, name)
    state.PushValue(-1)
    state.SetField(-2, "__index")
}

// Push a bufr message onto the stack and set its corresponding metatable
func pushMessage(state *lua.State, message *bufr.Message) {
    state.PushUserData(message)
    lua.SetMetaTableNamed(state, MESSAGE_META_TABLE)
}

// Push a section onto the stack and set its corresponding metatable
func pushSection(state *lua.State, section *bufr.Section) {
    state.PushUserData(section)
    lua.SetMetaTableNamed(state, SECTION_META_TABLE)
}

// Push a field onto the stack and set its corresponding metatable
func pushField(state *lua.State, field *bufr.Field) {
    state.PushUserData(field)
    lua.SetMetaTableNamed(state, FIELD_META_TABLE)
}

// Push the value of a field to stack
func pushFieldValue(state *lua.State, field *bufr.Field) error {
    switch field.Value.(type) {
    case *tdcfio.Binary:
        state.PushUserData(field.Value.(*tdcfio.Binary))
    default:
        return pushSimpleValue(state, field.Value)
    }
    return nil
}

// push a simple value onto the stack
func pushSimpleValue(state *lua.State, value interface{}) error {
    switch value.(type) {
    case int:
        state.PushInteger(value.(int))
    case uint:
        state.PushUnsigned(value.(uint))
    case string:
        state.PushString(value.(string))
    case bool:
        state.PushBoolean(value.(bool))
    case float64:
        state.PushNumber(value.(float64))
    default:
        return fmt.Errorf("simple value push not implemented for: %T", value)
    }

    return nil
}
