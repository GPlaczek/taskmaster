*** Settings ***
Resource  ./common.robot

*** Test Cases ***
GET Events End Point Should Work
    [Template]  Simple GET On Collection
    events

POST Events End Point Should Work
    [Template]  Simple POST On Collection
    events

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
