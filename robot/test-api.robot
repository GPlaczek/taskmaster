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

Add Test Event
    ${Body}  Create Dictionary
    Set To Dictionary  ${Body}  name  test
    Set To Dictionary  ${Body}  description  testtest
    Set To Dictionary  ${Body}  date  2024-05-06T08:40:00Z

    ${Response}  POST  http://localhost:8080/events  json=${body}

    ${Response Body}  Set Variable  ${Response.json()}
    ${Response Headers}  Set Variable  ${Response.headers}
    
    ${Event Id}  Set Variable  ${Response Body["id"]}
    ${ETag}  Set Variable  ${Response Headers["ETag"]}

    RETURN  ${Event Id}  ${ETag}


Get Test Event
    [Arguments]  ${id}
    ${Response}  GET  http://localhost:8080/events/${id}
    ${Response Body}  Set Variable  ${Response.json()}
    ${Response Headers}  Set Variable  ${Response.headers}
    
    ${ETag}  Set Variable  ${Response Headers["ETag"]}
    RETURN  ${Response Body}  ${ETag}
    

*** Test Cases ***
GET Events End Point Should Work
    ${response}  GET  http://localhost:8080/events

Adding Events Should Work
    ${Event Id}  ${_}  Add Test Event
    GET  http://localhost:8080/events/${Event Id}

Adding Events Should Not Work For Incomplete Data
    ${Event Data}  Create Dictionary
    Set To Dictionary  ${Event Data}  name  test
    Set To Dictionary  ${Event Data}  description  testtest

    POST  http://localhost:8080/events  json=${EventData}  expected_status=400

Updating Events Should Work
    ${Event Id}  ${ETag}  Add Test Event

    ${Request Body}  ${_}  Get Test Event  ${Event Id}
    Set To Dictionary  ${Request Body}  name  test1

    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  ${ETag}

    ${Response}  PUT  http://localhost:8080/events/${Event Id}  json=${Request Body}  headers=${Request Headers}

Update Should Fail For Missing ETag
    ${Event Id}  ${_}  Add Test Event

    ${GET Response}  GET  http://localhost:8080/events/${Event Id}

    ${Request Body}  Set Variable  ${GET Response.json()}
    Set To Dictionary  ${Request Body}  name  test1

    ${Response}  PUT  http://localhost:8080/events/${Event Id}  json=${Request Body}  expected_status=409

Update Should Fail For Invalid Job Id
    PUT  http://localhost:8080/events/99  expected_status=404

Update Should Fail For Incomplete Data
    ${Event Id}  ${_}  Add Test Event
    ${Event Data}  ${ETag}  Get Test Event  ${Event Id}

    Remove From Dictionary  ${Event Data}  name

    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  ${ETag}

    PUT  http://localhost:8080/events/${Event Id}  headers=${Request Headers}  json=${Event Data}  expected_status=400

Deleting Events Should Work
    ${Event Id}  ${ETag}  Add Test Event 
    
    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  ${ETag}

    DELETE  http://localhost:8080/events/${Event Id}  headers=${Request Headers}

Delete Should Fail For Invalid Job Id
    DELETE  http://localhost:8080/events/99  expected_status=404

Delete Should Fail For Missing ETag
    ${Event Id}  ${_}  Add Test Event
    
    DELETE  http://localhost:8080/events/${Event Id}  expected_status=409
