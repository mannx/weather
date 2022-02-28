# weather

Logs weather data using OpenWeatherMap.org into a database.


* Things to note:
	* Can store and view multiple cities
	* Early versions need to start with -city to update the database with initial city id's
	* Build script currently only builds for ARM with a 'beta' tag

* TODO:
	* Redo import code.  Parsing OWM response.
	* Update build.sh to either push or output a .tar.gz file
		- Should default to a .tar.gz unless specified
	* Set log levels via config file

* Installation:
	* See docker-compose below
	* Environment variables of note:
		* APIKEY  					->  api key for OpenWeatherMap.org
		* path : /data				->	path to where the database file will be located along with other app data (if required)
		* WEATHER_UPDATE_SCHEDULE	->	cron expression of when to pull new weather data in
		* WEATHER_DATA_PATH			->	if we want/need to change the data directory from /data to something else (Used mainly for development)

# Build Instructions

1. Run build.sh to build the docker image, and downnload and include the city list for openweathermap.org
2. Create a config.yml file in the data directory (see below)
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
       volumes:
            - *PATH TO DATA DIRECTORY*:/data
		restart: always
```

## Config.yml Example

```yaml
CityIDs: [id1, id2, ...]
APIKey: "API KEY HERE"
```
