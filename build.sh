#!/bin/sh

#
# Script is used to build the docker image for the weather app
# It downloads the city list from openweathermap.org if found
# and output the docker image in a .tar.gz ready for ARM/v7
# TODO: include output options for AMD also, along with --push flag if requested

# download and uncompress the city list in the current directory
getCityList() {
	url=http://bulk.openweathermap.org/sample/city.list.min.json.gz
	wget $url 

	if [ $? -ne 0 ]; then
		echo "Unable to download city list!"
		return 2
	fi

	# uncompress the file
	gzip -d city.list.min.json.gz
	if [ $? -ne 0 ]; then
		echo "Unable to decompress city list!"
		return 3
	fi

	return 1
}

args_help() {
	echo "$0 [-t TAG] [-T] [-l]"
	echo "-t TAG	=>	Set the tag to TAG"
	echo "-T		=>	Set the tag to the version number being built"
	echo "-l		=>	Build a local image without namespace, otherwise uses namespace and pushes to repo"
	return
}

echo "Building weather..."

# make sure we have the city list available as city.json
if [ ! -f "city.list.min.json" ]; then
	echo "city.json not found, downloading..."
	getCityList
	if [ ! $? -eq 1 ]; then
		echo "Build failed."
		exit 1
	fi
fi

# set default paramaters
TAG=dev
LOCAL=0

# read in build flags, determine if we are building a local dev image, or pushing a full image and grab the tag if provided
while getopts "t:lT" opt; do
	case $opt in
		t)
			TAG=$OPTARG
			echo "-t was providing, save tag..."
			;;
		l)
			LOCAL=1
			echo "-l was provided, building local image..."
			;;
		T)
			echo "-T provided, retrieving build tag from go-build.sh ..."
			TAG=$(./go-build.sh -t)
			;;
		\?)
			echo "Invalid option: -$OPTARG"
			args_help()
			exit 1
			;;
		:)
			echo "Option -$OPTARG requires an argument..."
			args_help()
			exit 1
			;;
	esac
done

echo "Building docker image with tag: $TAG..."

if [ $LOCAL -eq 1 ]; then
	echo "Building local image..."
	docker build -t weather:$TAG -t weather:latest .
else
	echo "Building repo image..."
	docker buildx build --platform linux/arm,linux/amd64 -t mannx/weather:$TAG -t mannx/weather:latest . --push
fi

echo "Docker build return with status $?"
