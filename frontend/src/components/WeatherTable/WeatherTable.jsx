import React from "react";
import WeatherChart from "./WeatherChart.jsx";
import UrlGet from "../URL/URL.jsx";


const Stats = ({High, Low, Rain, Snow}) => (
	<div><ul>
		<li>High: {High}</li>
		<li>Low: {Low}</li>
		<li>Snow: {Snow}</li>
		<li>Rain: {Rain}</li>
	</ul></div>
);

const Latest = ({Temp, FeelsLike, Wind}) => (
	<div><ul>
		<li>Temp: {Temp}</li>
		<li>Feels Like: {FeelsLike}</li>
		<li>Wind: {Wind}</li>
	</ul></div>
);

class WeatherTable extends React.Component {
	
	state = {
		loading: true,
		weather: null,

		cities: [], // list of cities we can see, defaults to the first one in the list
		cityid: 0, // city we are viewing

		error: false,
		errMsg: null,
	}	


	async componentDidMount() {
		await this.loadCities();	// make sure we have the data loaded before trying to use it

		// if we have citie id's, default to the first and display its data
		if(this.state.cities.length > 0) {
			const id = this.state.cities[0].ID;
			this.setState({cityid: id});
			this.loadData(id);
		}
	}

	loadData = async (id) => {
		const url = UrlGet("24hr") + "?id="+id;

		const resp = await fetch(url);
		const data = await resp.json();
		
		this.setState({loading: false, weather: data});
	}

	loadCities = async () => {
		const url = UrlGet("CityList");
		const resp = await fetch(url);
		const data = await resp.json();

		if(data === null) {
			// no data recieved
			this.setState({cities: [], error: true, errMsg: "No City Setup"})
		}else if(!data.Error){
			this.setState({cities: data});
		}else{
			this.setState({error: true, errMsg: "Unable to retrieve city list"});
		}
	}

	renderCities = () => {
		return (<div>
			<label>Pick a city to view
				<select value={this.state.cityid} onChange={(e) => this.cityUpdate(e)} >
					{this.state.cities.map(function(obj, i) {
						return <option key={obj.ID} value={obj.ID}>{obj.Name}</option>;
					})}
				</select>
			</label>
		</div>);
	}

	render() {
		let data = null;
		if(this.state.dataing || !this.state.weather) {
			data = <div>Loading current weather...</div>;
		}else{
			data = this.renderData();
		}

		return (<>
			{this.renderCities()}
			{data}
		</>);
	}

	cityUpdate = (e) => {
		this.setState({cityid: e.target.value});
		this.loadData(e.target.value);
	}

	renderData = () => {
		const size = this.state.weather.Data.length;
		if(size <= 0) { 
			return <h3>No data for the last 24 hours</h3>;
		}

		const wd = this.state.weather.Data[size-1];

		return (
			<div>
				<div>
					<h3>24 Hour Stats</h3>
					<Stats High={this.state.weather.High} Low={this.state.weather.Low} Snow={this.state.weather.Snow} Rain={this.state.weather.Rain} />
				</div>
				<div>
					<h3>Latest Weather</h3>
					<Latest Temp={wd.Temp} FeelsLike={wd.FeelsLike} Wind={wd.WindSpeed} />
				</div>
				<div>
					<span>Temperature</span>
					<WeatherChart data={this.state.weather.ChartData} item="Temp" />
				</div>
				<div>
					<span>Feels Like</span>
					<WeatherChart data={this.state.weather.ChartData} item="FeelsLike" />
				</div>

				<div><span>Snow</span>
					<WeatherChart data={this.state.weather.ChartData} item="Snow" />
				</div>

				<div><span>Rain</span>
					<WeatherChart data={this.state.weather.ChartData} item="Rain" />
				</div>

			</div>
		);
	}
}

export default WeatherTable;
