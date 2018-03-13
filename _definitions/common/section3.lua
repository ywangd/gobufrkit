local function deserialise()
    local section = factory.newSection(3, 'Data Description Section')

    local lengthInBytes = factory.newField {
        'lengthInBytes', UINT, 24;
    }

    factory.newField {
        'reservedBits', BIN, 8;
    }

    factory.newField {
        'nSubsets', UINT, 16;
        proxy = true,
    }

    factory.newField {
        'isObservation', BOOL, 1;
        proxy = true,
    }

    factory.newField {
        'isCompressed', BOOL, 1;
        proxy = true,
    }

    factory.newField {
        'flagBits', BIN, 6;
    }

    factory.newTemplateField {
        'unexpandedTemplate', 2, 6, 8, lengthInBytes:value();
    }

    factory.padding(lengthInBytes:value())

    return section
end

return {
    deserialise = deserialise
}