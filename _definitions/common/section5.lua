local function deserialise()
    local section = factory.newSection(5, 'End Section')

    factory.newField {
        'stopSignature', BYTES, 32;
        assert = function(x)
            assert(x == '7777', 'Stop signature not matched: ' .. x)
        end
    }

    return section
end

return {
    deserialise = deserialise
}