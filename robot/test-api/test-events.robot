*** Settings ***
Resource  ./common.robot

*** Test Cases ***
GET Events End Point Should Work
    [Template]  Simple GET On Collection
    events

POST Events End Point Should Work
    [Template]  Simple POST On Collection
    events

Adding Events Should Not Work For Incomplete Data
    ${Event Data}  Create Dictionary
    Set To Dictionary  ${Event Data}  name  test
    Set To Dictionary  ${Event Data}  description  testtest

    POST  http://localhost:8080/events  json=${EventData}  expected_status=400

Updating Events Should Work
    [Template]  Simple PUT On Object
    events

Updating Event Should Fail For Missing ETag
    [Template]  PUT Without ETag On Object
    events

Updating Event Should Fail For Invalid ETag
    [Template]  PUT With Invalid ETag On Object
    events

Updating Events Should Fail For Invalid Job Id
    [Template]  PUT With Invalid Id
    events

Updating Events Should Fail For Incomplete Data
    [Template]  PUT With Incomplete Data
    events

Deleting Events Should Work
    [Template]  Simple DELETE On Object
    events

Delete Should Fail For Invalid Event Id
    [Template]  DELETE With Invalid Id
    events

Deleting Events Should Fail For Missing ETag
    ${Event Id}  ${_}  Add Test Object  events
    DELETE  http://localhost:8080/events/${Event Id}  expected_status=428

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
