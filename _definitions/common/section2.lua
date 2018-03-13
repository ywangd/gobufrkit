local function deserialise()
    local section = factory.newSection(2, 'Optional Section')

    local lengthInBytes = factory.newField {
        'lengthInBytes', UINT, 24;
    }

    factory.newField {
        'reservedBits', BINARY, 8;
    }

    -- TODO: process local bits
    factory.newField {
        'localBits', BINARY, (lengthInBytes:value() - 4) * BITS_PER_BYTE;
    }

    return section
end

return {
    deserialise = deserialise
}