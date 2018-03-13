local function deserialise()
    local section = factory.newSection(4, 'Data Section')

    local lengthInBytes = factory.newField {
        'lengthInBytes', UINT, 24;
    }

    factory.newField {
        'reservedBits', BIN, 8;
    }

    factory.newPayloadField {
        'payload',
        factory.getMessage():getProxyField("nSubsets"):value(),
        factory.getMessage():getProxyField("isCompressed"):value(),
    }

    factory.padding(lengthInBytes:value())

    return section
end

return {
    deserialise = deserialise
}