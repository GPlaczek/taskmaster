*** Settings ***
Resource  ./common.robot

*** Test Cases ***
Merging Events Should Work
    ${Event Id 1}  ${_}  Add Test Object  events
    ${Event Id 2}  ${_}  Add Test Object  events

    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id1  ${Event Id 1}
    Set To Dictionary  ${Request Body}  id2  ${Event Id 2}

    POST  http://localhost:8080/merges  json=${Request Body}

Reading Merges Should Work
    ${Event Id 1}  ${_}  Add Test Object  events
    ${Event Id 2}  ${_}  Add Test Object  events

    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id1  ${Event Id 1}
    Set To Dictionary  ${Request Body}  id2  ${Event Id 2}

    POST  http://localhost:8080/merges  json=${Request Body}
    GET  http://localhost:8080/merges

Merging Event With Itself Should Not Work
    ${Event Id 1}  ${_}  Add Test Object  events

    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id1  ${Event Id 1}
    Set To Dictionary  ${Request Body}  id2  ${Event Id 1}

    POST  http://localhost:8080/merges  json=${Request Body}  expected_status=409

Merging Event With Invalid Event Should Not Work
    ${Event Id 1}  ${_}  Add Test Object  events

    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id1  ${Event Id 1}
    Set To Dictionary  ${Request Body}  id2  ${999999}

    POST  http://localhost:8080/merges  json=${Request Body}  expected_status=404

Merging Invalid Event With Event Should Not Work
    ${Event Id 1}  ${_}  Add Test Object  events

    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id1  ${999999}
    Set To Dictionary  ${Request Body}  id2  ${Event Id 1}

    POST  http://localhost:8080/merges  json=${Request Body}  expected_status=404

GET Should Work For New Event
    ${Event Id 1}  ${_}  Add Test Object  events
    ${Event Id 2}  ${_}  Add Test Object  events

    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id1  ${Event Id 1}
    Set To Dictionary  ${Request Body}  id2  ${Event Id 2}

    ${Response}  POST  http://localhost:8080/merges  json=${Request Body}
    ${Response Body}  Set Variable  ${Response.json()}
    ${Event Id}  Set Variable  ${Response Body["new_id"]}

    GET  http://localhost:8080/events/${Event Id}

GET Should Not Work For Old Events
    ${Event Id 1}  ${_}  Add Test Object  events
    ${Event Id 2}  ${_}  Add Test Object  events

    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id1  ${Event Id 1}
    Set To Dictionary  ${Request Body}  id2  ${Event Id 2}

    ${Response}  POST  http://localhost:8080/merges  json=${Request Body}
    ${Response Body}  Set Variable  ${Response.json()}

    GET  http://localhost:8080/events/${Event Id 1}  expected_status=404
    GET  http://localhost:8080/events/${Event Id 2}  expected_status=404

Bound Attachments Should Be Readable After Merging
    ${Event Id 1}  ${_}  Add Test Object  events
    ${Event Id 2}  ${_}  Add Test Object  events

    ${Attachment Id 1}  ${_}  Add Test Object  attachments
    ${Attachment Id 2}  ${_}  Add Test Object  attachments

    Bind Attachment To Event  ${Event Id 1}  ${Attachment Id 1}
    Bind Attachment To Event  ${Event Id 2}  ${Attachment Id 2}

    ${Request Body}  Create Dictionary
    Set To Dictionary  ${Request Body}  id1  ${Event Id 1}
    Set To Dictionary  ${Request Body}  id2  ${Event Id 2}

    ${Response}  POST  http://localhost:8080/merges  json=${Request Body}
    ${Response Body}  Set Variable  ${Response.json()}
    ${Event Id}  Set Variable  ${Response Body["new_id"]}

    GET  http://localhost:8080/events/${Event Id}/attachments
    GET  http://localhost:8080/events/${Event Id}/attachments/${Attachment Id 1}
    GET  http://localhost:8080/events/${Event Id}/attachments/${Attachment Id 2}
