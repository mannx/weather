# weather

Logs weather data using OpenWeatherMap.org into a database.

* Currently Working On:
	* Add/remove cities to monitor from the settings page
	* Implement removal of cities and include state/country when listing current ones

* Things to note:
	* Can store and view multiple cities
	* Early versions need to start with -city to update the database with initial city id's
	* Build script currently only builds for ARM with a 'beta' tag

* TODO:
	* Redo import code.  Parsing OWM response.
	* Update build.sh to either push or output a .tar.gz file
		- Should default to a .tar.gz unless specified
	* [PARTIAL] Adding/Removing watch cities is now part of the settings web page

* Installation:
	* See docker-compose below
	* Environment variables of note:
		* path : /data				->	path to where the database file will be located along with other app data
		* WEATHER_UPDATE_SCHEDULE	->	cron expression of when to pull new weather data in
		* WEATHER_DATA_PATH			->	if we want/need to change the data directory from /data to something else (Used mainly for development)
		* WEATHER_API_KEY			->  api key for OpenWeatherMap.org
		* WEATHER_LOG_LEVEL			->	[OPTIONAL] Level for logging output, defaults to INFO (see below)

### Log Levels

The following are options for setting log levels
* DEBUG
* INFO
* WARN
* ERROR
* FATAL
* PANIC
* NOLEVEL
* DISABLED

Note: DISABLED prevents logger output all together

# Build Instructions

1. Run build.sh to build the docker image, and downnload and include the city list for openweathermap.org
2. ~~Create a config.yml file in the data directory (see below)~~
3. Start the container either on the commandline or with a docker-compose file (see below for example)

## Docker Compose Example

```dockerfile
version: "2.0"

services:
	weather:
    	image: mannx/weather
        container_name: weather
        ports:
			- 8080:8080
        environment:
            - TZ=*TIMEZONE HERE*
            - WEATHER_UPDATE_SCHEDULE=*/30 * * * *
			- WEATHER_API_KEY= --API KEY HERE--
       volumes:
            - *PATH TO DATA DIRECTORY*:/data
		restart: always
```
