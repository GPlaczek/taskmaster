*** Settings ***
Resource  ./common.robot

*** Test Cases ***
GET Attachments End Point Should Work
    [Template]  Simple GET On Collection
    attachments

POST Attachments End Point Should Work
    [Template]  Simple POST On Collection
    attachments

Updating Attachments Should Work
    [Template]  Simple PUT On Object
    attachments

Updating Attachment Should Fail For Missing ETag
    [Template]  PUT Without ETag On Object
    attachments

Updating Attachment Should Fail For Invalid ETag
    [Template]  PUT With Invalid ETag On Object
    attachments

Updating Attachments Should Fail For Invalid Job Id
    [Template]  PUT With Invalid Id
    attachments

Updating Attachments Should Fail For Incomplete Data
    [Template]  PUT With Incomplete Data
    attachments

Deleting Attachments Should Work
    [Template]  Simple DELETE On Object
    attachments

Delete Should Fail For Invalid Attachment Id
    [Template]  DELETE With Invalid Id
    attachments

Deleting Attachments Should Fail For Missing ETag
    ${Attachment Id}  ${_}  Add Test Object  attachments
    DELETE  http://localhost:8080/attachments/${Attachment Id}  expected_status=428
