
const production = process.env.NODE_ENV !== 'development';
const baseURL = "http://localhost:8080";

const urls = {
	"24hr": "/api/24hr",
	"Latest": "/api/latest",
	"Daily": "/api/daily",
	"Weekly": "/api/weekly",
	"CityList": "/api/cities",
}

function UrlGet(name) {
	var base = "";
	if(!production) {
		base = baseURL;
	}

	return base+urls[name];
}

export default UrlGet;
