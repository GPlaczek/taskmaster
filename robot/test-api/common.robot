*** Settings ***
Library   Collections
Library   RequestsLibrary

*** Keywords ***
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
