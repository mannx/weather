# weather

Logs weather data using OpenWeatherMap.org into a database.

- Things to note:

- TODO:
	* Better logging
	* Output needs to be done
	* 

- Installation:
	* See docker-compose below
	* Environment variables of note:
		* APIKEY  					->  api key for OpenWeatherMap.org
		* path : /data				->	path to where the database file will be located along with other app data (if required)
		* WEATHER_UPDATE_SCHEDULE	->	cron expression of when to pull new weather data in
		* WEATHER_DATA_PATH			->	if we want/need to change the data directory from /data to something else (Used mainly for development)


## Docker Compose Example

```dockerfile
version: "3.0"

services:
        weather:
                image: mannx/weather
                container_name: weather
                ports:
                        - 8080:8080
                environment:
                        - APIKEY=*API KEY HERE*
                        - TZ=*TIMEZONE HERE*
                        - WEATHER_UPDATE_SCHEDULE=*/30 * * * *
                volumes:
                        - *PATH TO DATA DIRECTORY*:/data
                deploy:
                        restart_policy:
                                condition: on-failure
```
