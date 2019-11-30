#!/bin/bash

DAY=$1
paddedday="$(printf "%02d" ${DAY})"

sed \
 -e "s/!DAY!/${paddedday}/g" \
 -e "s/MAIN/main/" \
template.go > day${paddedday}.go
