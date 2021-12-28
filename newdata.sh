#!/bin/sh

#
# This script is run on a regular schedule to
# pull in new weather data to the database
# This method should be changed, atleast the method used to do so
#

curl localhost:8080/new
