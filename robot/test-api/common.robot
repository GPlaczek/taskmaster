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
    ${POST Response}  POST  http://localhost:8080/${collection}  expected_status=201
    ${POST Response Headers}  Set Variable  ${POST Response.headers}
    ${Location}  Set Variable  ${POST Response Headers["Location"]}
    ${ETag}  Set Variable  ${POST Response Headers["ETag"]}

    ${PUT Request Body}  Create Dictionary
    IF  '${collection}' == 'events'
        Set To Dictionary  ${PUT Request Body}  name  test
        Set To Dictionary  ${PUT Request Body}  description  testtest
        Set To Dictionary  ${PUT Request Body}  date  2024-05-06T08:40:00Z
    ELSE
        Set To Dictionary  ${PUT Request Body}  data  Test content
    END
    ${PUT Request Headers}  Create Dictionary
    Set To Dictionary  ${PUT Request Headers}  If-Match  ${ETag}

    ${PUT Response}  PUT  http://localhost:8080${Location}  json=${PUT Request Body}  headers=${PUT Request Headers}
    ${PUT Response Body}  Set Variable  ${PUT Response.json()}
    ${Id}  Set Variable  ${PUT Response Body["id"]}

    RETURN  ${Id}  ${ETag}

Bind Attachment To Event
    [Arguments]  ${eid}  ${aid}  ${expected_status}=201
    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id  ${aid}

    POST  http://localhost:8080/events/${eid}/attachments  json=${Request Body}  expected_status=${expected_status}

Unbind Attachment From Event
    [Arguments]  ${eid}  ${aid}  ${expected_status}=204

    DELETE  http://localhost:8080/events/${eid}/attachments/${aid}  expected_status=${expected_status}

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

    ${Response}  PUT  http://localhost:8080/${collection}/${Id}  json=${Request Body}  expected_status=428

PUT With Invalid ETag On Object
    [Arguments]  ${collection}

    ${Id}  ${_}  Add Test Object  ${collection}
    ${Response}  GET  http://localhost:8080/${collection}/${Id}
    ${Response Body}  Set Variable  ${Response.json()}
    ${Request Body}  Get Updated Object Body  ${collection}  ${Response Body}

    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  deadbeefcafe12345678

    ${Response}  PUT  http://localhost:8080/${collection}/${Id}  json=${Request Body}  expected_status=412  headers=${Request Headers}

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

    DELETE  http://localhost:8080/${collection}/${Id}  headers=${Request Headers}  expected_status=204

PUT With Invalid Id
    [Arguments]  ${collection}
    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  deadbeefcafe12345678

    PUT  http://localhost:8080/${collection}/999999  expected_status=404  headers=${Request Headers}  data={}

DELETE With Invalid Id
    [Arguments]  ${collection}
    ${Request Headers}  Create Dictionary
    Set To Dictionary  ${Request Headers}  If-Match  deadbeefcafe12345678

    DELETE  http://localhost:8080/${collection}/999999  expected_status=404  headers=${Request Headers}
    
