*** Settings ***
Resource  ./common.robot

*** Test Cases ***
GET Attachments End Point Should Work
    [Template]  Simple GET On Collection
    attachments

POST Attachments End Point Should Work
    [Template]  Simple POST On Collection
    attachments

Adding Attachments Should Not Work For Incomplete Data
    ${Event Data}  Create Dictionary

    POST  http://localhost:8080/attachments  json=${EventData}  expected_status=400

Updating Attachments Should Work
    [Template]  Simple PUT On Object
    attachments

Updating Attachment Should Fail For Missing ETag
    [Template]  PUT Without ETag On Object
    attachments

Updating Attachments Should Fail For Invalid Job Id
    PUT  http://localhost:8080/attachments/-1  expected_status=404  data={}

Updating Attachments Should Fail For Incomplete Data
    [Template]  PUT With Incomplete Data
    attachments

Deleting Attachments Should Work
    [Template]  Simple DELETE On Object
    attachments

Delete Should Fail For Invalid Attachment Id
    DELETE  http://localhost:8080/attachments/-1  expected_status=404

Deleting Attachments Should Fail For Missing ETag
    ${Attachment Id}  ${_}  Add Test Object  attachments
    DELETE  http://localhost:8080/attachments/${Attachment Id}  expected_status=409
