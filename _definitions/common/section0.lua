local function deserialise()
    local section = factory.newSection(0, 'Indicator Section')

    factory.newField {
        'startSignature', BYTES, 32;
    }

    factory.newField {
        'totalLengthInBytes', UINT, 24;
        proxy = true,
    }

    factory.newField {
        'bufrEditionNumber', UINT, 8;
        proxy = true,
    }

    return section
end

return {
    deserialise = deserialise
}