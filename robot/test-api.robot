*** Settings ***
Library        Process
Library        RequestsLibrary
Library        Collections

Test Timeout   5s
Test Setup     Setup Test
Test Teardown  Teardown Test

*** Keywords ***
Setup Test
    ${Server Process}=   Start Process  ../server  stdout=log.txt
    Sleep                0.1s
    Set Global Variable  ${Server}      ${Server Process}

Teardown Test
    Terminate Process    ${Server}

Get Test Object
    [Arguments]  ${collection}  ${id}
    ${Response}  GET  http://localhost:8080/${collection}/${id}
    ${Response Body}  Set Variable  ${Response.json()}
    ${Response Headers}  Set Variable  ${Response.headers}
   
    ${ETag}  Set Variable  ${Response Headers["ETag"]}
    RETURN  ${Response Body}  ${ETag}

Add Test Object
    [Arguments]  ${collection}
    ${Body}  Create Dictionary
    IF  '${collection}' == 'events'
        Set To Dictionary  ${Body}  name  test
        Set To Dictionary  ${Body}  description  testtest
        Set To Dictionary  ${Body}  date  2024-05-06T08:40:00Z
    ELSE
        Set To Dictionary  ${Body}  data  Test content
    END

    ${Response}  POST  http://localhost:8080/${collection}  json=${body}

    ${Response Body}  Set Variable  ${Response.json()}
    ${Response Headers}  Set Variable  ${Response.headers}

    ${Id}  Set Variable  ${Response Body["id"]}
    ${ETag}  Set Variable  ${Response Headers["ETag"]}

    RETURN  ${Id}  ${ETag}

Bind Attachment To Event
    [Arguments]  ${eid}  ${aid}  ${expected_status}=200
    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id  ${aid}

    POST  http://localhost:8080/events/${eid}/attachments  json=${Request Body}  expected_status=${expected_status}

Simple GET On Collection
    [Arguments]  ${collection}
    ${response}  GET  http://localhost:8080/${collection}

Simple POST On Collection
    [Arguments]  ${collection}
    ${Id}  ${_}  Add Test Object  ${collection}
    GET  http://localhost:8080/${collection}/${Id}

Get Updated Object Body
    [Arguments]  ${collection}  ${body}
   
    IF  '${collection}' == 'events'
        Set To Dictionary  ${body}  name  test1
    ELSE
        Set To Dictionary  ${body}  data  Test1 content
    END

    RETURN  ${body}

Simple PUT On Object
    [Arguments]  ${collection}

    ${Id}  ${ETag}  Add Test Object  ${collection}
    ${Response}  GET  http://localhost:8080/${collection}/${Id}
    ${Response Body}  Set Variable  ${Response.json()}
    ${Request Body}  Get Updated Object Body  ${collection}  ${Response Body}

    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  ${ETag}

    ${Response}  PUT  http://localhost:8080/${collection}/${Id}  json=${Request Body}  headers=${Request Headers}

PUT Without ETag On Object
    [Arguments]  ${collection}

    ${Id}  ${_}  Add Test Object  ${collection}
    ${Response}  GET  http://localhost:8080/${collection}/${Id}
    ${Response Body}  Set Variable  ${Response.json()}
    ${Request Body}  Get Updated Object Body  ${collection}  ${Response Body}

    ${Response}  PUT  http://localhost:8080/${collection}/${Id}  json=${Request Body}  expected_status=409

PUT With Incomplete Data
    [Arguments]  ${collection}
    ${Id}  ${_}  Add Test Object  ${collection}
    ${Event Data}  ${ETag}  Get Test Object  ${collection}  ${Id}

    IF  '${collection}' == 'events'
        Remove From Dictionary  ${Event Data}  name
    ELSE
        Remove From Dictionary  ${Event Data}  data
    END

    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  ${ETag}

    PUT  http://localhost:8080/${collection}/${Id}  headers=${Request Headers}  json=${Event Data}  expected_status=400

Simple DELETE On Object
    [Arguments]  ${collection}
    ${Id}  ${ETag}  Add Test Object  ${collection}

    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  ${ETag}

    DELETE  http://localhost:8080/${collection}/${Id}  headers=${Request Headers}


*** Test Cases ***
GET Events End Point Should Work
    [Template]  Simple GET On Collection
    events

GET Attachments End Point Should Work
    [Template]  Simple GET On Collection
    attachments

POST Events End Point Should Work
    [Template]  Simple POST On Collection
    events

POST Attachments End Point Should Work
    [Template]  Simple POST On Collection
    attachments

Adding Events Should Not Work For Incomplete Data
    ${Event Data}  Create Dictionary
    Set To Dictionary  ${Event Data}  name  test
    Set To Dictionary  ${Event Data}  description  testtest

    POST  http://localhost:8080/events  json=${EventData}  expected_status=400

Adding Attachments Should Not Work For Incomplete Data
    ${Event Data}  Create Dictionary

    POST  http://localhost:8080/attachments  json=${EventData}  expected_status=400

Updating Events Should Work
    [Template]  Simple PUT On Object
    events

Updating Attachments Should Work
    [Template]  Simple PUT On Object
    attachments

Updating Event Should Fail For Missing ETag
    [Template]  PUT Without ETag On Object
    events

Updating Attachment Should Fail For Missing ETag
    [Template]  PUT Without ETag On Object
    attachments

Updating Events Should Fail For Invalid Job Id
    PUT  http://localhost:8080/events/-1  expected_status=404  data={}

Updating Attachments Should Fail For Invalid Job Id
    PUT  http://localhost:8080/attachments/-1  expected_status=404  data={}

Updating Events Should Fail For Incomplete Data
    [Template]  PUT With Incomplete Data
    events

Updating Attachments Should Fail For Incomplete Data
    [Template]  PUT With Incomplete Data
    attachments

Deleting Events Should Work
    [Template]  Simple DELETE On Object
    events

Deleting Attachments Should Work
    [Template]  Simple DELETE On Object
    attachments

Delete Should Fail For Invalid Event Id
    DELETE  http://localhost:8080/events/-1  expected_status=404

Delete Should Fail For Invalid Attachment Id
    DELETE  http://localhost:8080/attachments/-1  expected_status=404

Deleting Events Should Fail For Missing ETag
    ${Event Id}  ${_}  Add Test Object  events
    DELETE  http://localhost:8080/events/${Event Id}  expected_status=409

Deleting Attachments Should Fail For Missing ETag
    ${Attachment Id}  ${_}  Add Test Object  attachments
    DELETE  http://localhost:8080/attachments/${Attachment Id}  expected_status=409

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

    Bind Attachment To Event  -1  ${Attachment Id}  expected_status=404

Binding Attachment Should Not Work For Invalid Attachment
    ${Event Id}  ${_}  Add Test Object  events
    ${Attachment Id}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id}  ${-1}  expected_status=404
