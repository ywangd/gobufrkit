package api

import (
    "github.com/Shopify/go-lua"
    "fmt"
    "github.com/ywangd/gobufrkit/bufr"
)

const MESSAGE_META_TABLE = "gobufrkit.message"

type ObjMessage struct {
}

// TODO: Improve the string format
func (obj *ObjMessage) messageToString(state *lua.State) int {
    message := state.ToUserData(1).(*bufr.Message)
    state.PushString(fmt.Sprintf("%+v", message))
    return 1
}

func (obj *ObjMessage) setProxyField(state *lua.State) int {
    message := state.ToUserData(1).(*bufr.Message)
    field := state.ToUserData(2).(*bufr.Field)
    message.SetProxyField(field)
    return 0
}

func (obj *ObjMessage) getProxyField(state *lua.State) int {
    message := state.ToUserData(1).(*bufr.Message)
    name, _ := state.ToString(2)
    field, err := message.ProxyField(name)
    if err != nil {
        state.PushString(err.Error())
        state.Error()
        return 0
    }
    pushField(state, field)

    return 1
}

func (obj *ObjMessage) registerGribMessageType(state *lua.State) {

    lua.NewMetaTable(state, MESSAGE_META_TABLE)

    // metatable.__index = metatable
    state.PushValue(-1)
    state.SetField(-2, "__index")

    lua.SetFunctions(state, []lua.RegistryFunction{
        {Name: "__tostring", Function: obj.messageToString},
        {Name: "setProxyField", Function: obj.setProxyField},
        {Name: "getProxyField", Function: obj.getProxyField},
    }, 0)

}
