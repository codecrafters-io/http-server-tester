#!/bin/bash

for i in $(seq 1 100)
do
    echo "Running iteration $i"
    make test_persistence_pass > /tmp/test
    if [ $? -ne 0 ]; then
        echo "make test_persistence_pass failed on iteration $i"
        exit 1
    fi
done    