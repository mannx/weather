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
	}	


	async componentDidMount() {
		const url = UrlGet("24hr");

		const resp = await fetch(url);
		const data = await resp.json();
		
		this.setState({loading: false, weather: data});
	}

	render() {
		if(this.state.loading || !this.state.weather) {
				return <div>Loading current weather...</div>;
		}

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
