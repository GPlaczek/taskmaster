*** Settings ***
Library        Process
Library        Collections

Suite Setup          Setup Test
Suite Teardown       Teardown Test

*** Keywords ***
Setup Test
    ${Server Process}=   Start Process  ../server  stdout=log.txt  stderr=log.txt
    Sleep                0.1s
    Set Global Variable  ${Server}      ${Server Process}

Teardown Test
    Terminate Process    ${Server}
