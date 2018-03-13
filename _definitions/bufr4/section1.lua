local function deserialise()
    local section = factory.newSection(1, 'Identification Section')

    local lengthInBytes = factory.newField {
        'lengthInBytes', UINT, 24;
    }

    factory.newField {
        'masterTableNumber', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'originatingCentre', UINT, 16;
        proxy = true,
    }

    factory.newField {
        'originatingSubCentre', UINT, 16;
        proxy = true,
    }

    factory.newField {
        'updateSequenceNumber', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'isSection2Presents', BOOL, 1;
        proxy = true,
    }

    factory.newField {
        'flagBits', BINARY, 7;
    }

    factory.newField {
        'dataCategory', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'dataI18nSubCategory', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'dataLocalSubCategory', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'masterTableVersion', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'localTableVersion', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'year', UINT, 16;
        proxy = true,
    }

    factory.newField {
        'month', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'day', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'hour', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'minute', UINT, 8;
        proxy = true,
    }

    factory.newField {
        'second', UINT, 8;
        proxy = true,
    }

    factory.padding(lengthInBytes:value())

    return section
end

return {
    deserialise = deserialise
}