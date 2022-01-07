# weather

Logs weather data using OpenWeatherMap.org into a database.

- Things to note:

- TODO:
	* Better logging
	* Output needs to be done

- Installation:
	* See docker-compose below
	* Environment variables of note:
		* APIKEY  					->  api key for OpenWeatherMap.org
		* path : /data				->	path to where the database file will be located along with other app data (if required)
		* WEATHER_UPDATE_SCHEDULE	->	cron expression of when to pull new weather data in
