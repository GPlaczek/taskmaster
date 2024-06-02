*** Settings ***
Resource  ./common.robot

*** Test Cases ***
Binding Attachments To Events Should Work
    ${Event Id}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id}  ${Attachment Id}

Reading Bound Attachments Should Work
    ${Event Id}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id}  ${Attachment Id}

    ${Response}  GET  http://localhost:8080/events/${Event Id}/attachments
    ${Body}  Set Variable  ${Response.json()}
    Should Not Be Empty  ${Body}

Binding The Same Attachment Twice Should Not Work
    ${Event Id}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id}  ${Attachment Id}
    Bind Attachment To Event  ${Event Id}  ${Attachment Id}  expected_status=409

Binding Attachment Should Not Work For Invalid Event
    ${Event Id}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  999999  ${Attachment Id}  expected_status=404

Binding Attachment Should Not Work For Invalid Attachment
    ${Event Id}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id}  ${999999}  expected_status=404

Binding Attachment That Is Already Bound Should Not Work
    ${Event Id 1}  ${_}  Add Test Object  events
    ${Event Id 2}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id 1}  ${Attachment Id}
    Bind Attachment To Event  ${Event Id 2}  ${Attachment Id}  expected_status=409

Unbinding Attachment Should Work
    ${Event Id}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id}  ${Attachment Id}
    Unbind Attachment From Event  ${Event Id}  ${Attachment Id}

    GET  http://localhost:8080/events/${Event Id}/attachments/${Attachment Id}  expected_status=404

Unbound Attachment Should Be Readable
    ${Event Id}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id}  ${Attachment Id}
    Unbind Attachment From Event  ${Event Id}  ${Attachment Id}

    GET  http://localhost:8080/attachments/${Attachment Id}  expected_status=200
