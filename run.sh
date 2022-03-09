#!/bin/sh

# run with the -city flag if the CITY file is not found in the data directory
flag=
if [ ! -f "./data/CITY" ]; then
	flag=-city
	touch ./data/CITY
fi

echo "Starting weather with flags: $flag..."
./weather $flag
