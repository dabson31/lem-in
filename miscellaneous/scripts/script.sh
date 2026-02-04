#!/bin/bash

for file in ../examples/*; do
    echo "Running with $file"
    go run ../main.go "$file"
    echo "---------------"

    read -p "Press ENTER to run next file or (Ctrl-c to stop)...."
    clear
done
