#!/bin/bash

set -eu

SCRIPT_LOCATION="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

function random_number() {
    local min=$1
    local max=$2

    echo $(( ( RANDOM % $max )  + $min ))
}

function count_lines() {
    local file=$1

    echo $(wc -l $file | awk '{print $1}')
}

function get_line() {
    local file=$1
    local line_number=$2

    echo $(sed -n "${line_number}p" $file)
}

function get_random_line() {
    local file=$1
    local line_number=$(random_number 1 $(count_lines $file))

    echo $(get_line $file $line_number)
}

function generate_envs() {
    local envs_file=$(find $PWD -type f -name "*.env.tpl")

    for env_file in $envs_file; do
        local env_name=$(basename $env_file .tpl)
        
        cat $env_file | envsubst > "$SCRIPT_LOCATION/$env_name"
    done
}

function main() {
    local users_file=$1
    local passwords_file=$2

    local username=$(get_random_line $users_file)
    local password=$(get_random_line $passwords_file)

    echo "Username: $username"
    echo "Password: $password"

    export USERNAME=$username
    export PASSWORD=$password

    generate_envs

    cd $SCRIPT_LOCATION
    docker compose up
}

if [ $# -ne 2 ]; then
    echo "Usage: $0 <users_file> <passwords_file>"
    exit 1
fi

USERSFILE=$1
PASSWORDSFILE=$2

main $USERSFILE $PASSWORDSFILE