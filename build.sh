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

# build the container and output in a .tar.gz
# NOTE: include option to pick a tag, currently defaults to beta

#docker buildx build --platform linux/arm -t mannx/weather:beta -o type=tar,dest=weather.tar .
#docker build -t weather .
docker buildx build --platform linux/arm,linux/amd64 -t mannx/weather:latest . --push
echo "Docker build return with status $?"

# if successfuly build, run tar through gz to get compressed archive
if [ $? -eq 0 ]; then
	echo "Build completed successfully..."
	echo "Compressing output archive..."
	gzip weather.tar
   	if [ $? -eq 1 ]; then
		echo "Compression successful"
	fi
else
	echo "Build Failed"
fi
