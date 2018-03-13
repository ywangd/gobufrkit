package api

import (
    "github.com/Shopify/go-lua"
    "github.com/ywangd/gobufrkit/deserialize"
)

const DESERIALIZER_META_TABLE = "gobufrkit.deserializer"

type LibDeserializer struct {
    factory deserialize.Factory
}

func (lib *LibDeserializer) getMessage(state *lua.State) int {
    pushMessage(state, lib.factory.Message())
    return 1
}

func (lib *LibDeserializer) newMessage(state *lua.State) int {
    inputPath, _ := state.ToString(1)
    message := lib.factory.NewMessage(inputPath)
    pushMessage(state, message)
    return 1
}

func (lib *LibDeserializer) newSection(state *lua.State) int {
    number, _ := state.ToInteger(1)
    description, _ := state.ToString(2)

    section := lib.factory.NewSection(number, description)
    pushSection(state, section)

    return 1
}

func (lib *LibDeserializer) initTableGroup(state *lua.State) int {
    masterTableNo, _ := state.ToInteger(1)
    centreNo, _ := state.ToInteger(2)
    subCentreNo, _ := state.ToInteger(3)
    wmoVersion, _ := state.ToInteger(4)
    localVersion, _ := state.ToInteger(5)
    err := lib.factory.InitTableGroup(masterTableNo, centreNo, subCentreNo, wmoVersion, localVersion)
    if err != nil {
        state.PushString(err.Error())
        state.Error()
        return 0
    }
    return 0
}

func (lib *LibDeserializer) newField(state *lua.State) int {

    state.RawGetInt(1, 1)
    name, _ := state.ToString(-1)

    state.RawGetInt(1, 2)
    dn, _ := state.ToInteger(-1)
    dataType := deserialize.DataType(dn)

    state.RawGetInt(1, 3)
    nbits, _ := state.ToInteger(-1)

    state.Field(1, "proxy")
    proxy := state.ToBoolean(-1)

    field, err := lib.factory.NewField(name, dataType, nbits, proxy)
    if err != nil {
        state.PushString(err.Error())
        state.Error()
        return 0
    }

    pushField(state, field)
    return 1
}

func (lib *LibDeserializer) newTemplateField(state *lua.State) int {
    state.RawGetInt(1, 1)
    name, _ := state.ToString(-1)

    state.RawGetInt(1, 2)
    fbits, _ := state.ToInteger(-1)

    state.RawGetInt(1, 3)
    xbits, _ := state.ToInteger(-1)

    state.RawGetInt(1, 4)
    ybits, _ := state.ToInteger(-1)

    state.RawGetInt(1, 5)
    sectionLengthInBytes, _ := state.ToUnsigned(-1)

    field, err := lib.factory.NewTemplateField(name, fbits, xbits, ybits, sectionLengthInBytes)
    if err != nil {
        state.PushString(err.Error())
        state.Error()
        return 0
    }
    pushField(state, field)
    return 1
}

func (lib *LibDeserializer) newPayloadField(state *lua.State) int {
    state.RawGetInt(1, 1)
    name, _ := state.ToString(-1)

    state.RawGetInt(1, 2)
    nsubsets, _ := state.ToInteger(-1)

    state.RawGetInt(1, 3)
    compressed := state.ToBoolean(-1)

    field, err := lib.factory.NewPayloadField(name, nsubsets, compressed)
    if err != nil {
        state.PushString(err.Error())
        state.Error()
        return 0
    }
    pushField(state, field)
    return 1
}

func (lib *LibDeserializer) padding(state *lua.State) int {
    sectionLengthInBytes, _ := state.ToUnsigned(1)
    _, err := lib.factory.Padding(sectionLengthInBytes)
    if err != nil {
        state.PushString(err.Error())
        state.Error()
        return 0
    }
    return 0
}

func (lib *LibDeserializer) peekEditionNumber(state *lua.State) int {
    v, err := lib.factory.PeekEditionNumber()
    if err != nil {
        state.PushString(err.Error())
        state.Error()
        return 0
    }
    state.PushUnsigned(v)
    return 1
}

func (lib *LibDeserializer) Register(state *lua.State) int {
    lua.NewLibrary(state, []lua.RegistryFunction{
        {Name: "getMessage", Function: lib.getMessage},
        {Name: "newMessage", Function: lib.newMessage},
        {Name: "newSection", Function: lib.newSection},
        {Name: "newField", Function: lib.newField},
        {Name: "newTemplateField", Function: lib.newTemplateField},
        {Name: "newPayloadField", Function: lib.newPayloadField},
        {Name: "padding", Function: lib.padding},
        {Name: "peekEditionNumber", Function: lib.peekEditionNumber},
        {Name: "initTableGroup", Function: lib.initTableGroup},
    })
    return 1
}
