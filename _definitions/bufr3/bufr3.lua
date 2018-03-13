local function deserialise()
    local section0 = require 'common.section0'
    local section1 = require 'bufr3.section1'
    local section2 = require 'common.section2'
    local section3 = require 'common.section3'
    local section4 = require 'common.section4'
    local section5 = require 'common.section5'

    local message = factory.newMessage()

    section0.deserialise()
    section1.deserialise()

    if message:getProxyField('isSection2Presents'):value() then
        section2.deserialise()
    end

    section3.deserialise()

    -- Config tables for lookup
    factory.initTableGroup(
        message:getProxyField('masterTableNumber'):value(),
        message:getProxyField('originatingCentre'):value(),
        message:getProxyField('originatingSubCentre'):value(),
        message:getProxyField('masterTableVersion'):value(),
        message:getProxyField('localTableVersion'):value()
    )

    section4.deserialise()
    section5.deserialise()

    return message
end

return {
    deserialise = deserialise
}