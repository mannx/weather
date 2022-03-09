
const production = process.env.NODE_ENV !== 'development';
const baseURL = "http://localhost:8080";

const urls = {
	"24hr": "/api/24hr",
	"Latest": "/api/latest",
	"Daily": "/api/daily",
	"Weekly": "/api/weekly",

	"CityList": "/api/cities",			// retrieve all current cities in config
	"CityAdd": "/api/city/add",			// add a new city, return validation information
	"CityConfirm": "/api/city/confirm",	// confirm the new city to be added

	"Migrate": "/api/migrate",
}

function UrlGet(name) {
	var base = "";
	if(!production) {
		base = baseURL;
	}

	return base+urls[name];
}

export default UrlGet;
