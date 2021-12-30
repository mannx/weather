# weather

Logs weather data using OpenWeatherMap.org into a database.

- Things to note:
	* update time is currently restricted to 15min, 1hour, or daily due to alpine limits
		~ Could be changed with different distro? size?
		+ No need to deal with cron jobs/expressions to get working, but limits flexibility
	* Container is started using a script to run the cron daemon, then launch the weather app

- TODO:
	* Better logging
	* Output needs to be done
	* Retrieving new weather is done with a GET using a cron job.  Find better way? 
	* Include some sort of version identification to confirm docker images are up to date

- Installation:
	* See docker-compose below
	* Environment variables of note:
		* APIKEY  		->  api key for OpenWeatherMap.org
		* path : /app	->	path to where the database file will be located along with other app data (if required)
	* To automate weather retrieval, the following cron job will need to be set up on the host 
		* @hourly docker exec /newdata.sh
