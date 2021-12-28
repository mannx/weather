#!/bin/sh

#
#	This script is used to start the container
#	It is used to make sure the cron daemon is running and then
#	executes the main program

rc-service crond start && rc-update add crond
/weather
