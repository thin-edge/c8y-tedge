*** Settings ***
Resource        ../resources/common.robot
Library         Cumulocity
Library        Process

Suite Setup     Set Main Device


*** Test Cases ***
Bootstrap main device to the cloud
    ${result}=    Process.Run Process    docker    compose    exec    c8y    zsh    -c    c8y tedge bootstrap root@tedge '${DEVICE_ID}'    cwd=${CURDIR}/../..    shell=True    timeout=1min
    Cumulocity.Device Should Exist    ${DEVICE_ID}

Enable mtls on main device
    Process.Run Process    docker compose exec c8y zsh -c "c8y tedge bootstrap root@tedge --mtls"

Enable mtls on child device
    Process.Run Process    docker compose exec c8y zsh -c "c8y tedge bootstrap root@child01 --mtls --main-device root@tedge"
