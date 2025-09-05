#!/bin/sh

sed_flag=''

# Determine the correct `sed` flag for in-place editing
if [ "$OSTYPE" = "darwin"* ]; then
    sed_flag='-i ""'
else
    sed_flag='-i'
fi

# Function to run sed commands
replace_in_files() {
    eval sed $sed_flag "s/x-nullable/nullable/g" "$1"
    eval sed $sed_flag "s/x-omitempty/omitempty/g" "$1"
    eval sed $sed_flag "s/x-example/example/g" "$1"
}

# Apply replacements to the files
replace_in_files "./docs/docs.go"
replace_in_files "./docs/swagger.json"
replace_in_files "./docs/swagger.yaml"
