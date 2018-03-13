--config:debug(true)

local editionNumber = factory.peekEditionNumber()

if editionNumber == 4 then
    local bufr4 = require 'bufr4.bufr4'
    return bufr4.deserialise()

elseif editionNumber == 3 then
    local bufr3 = require 'bufr3.bufr3'
    return bufr3.deserialise()

else
    error("Invalid BUFR edition number: " .. editionNumber)
end