import React from "react";
import WeatherChart from "./WeatherChart.jsx";


const Stats = ({High, Low, Rain, Snow}) => (
		<div><ul>
				<li>High: {High}</li>
				<li>Low: {Low}</li>
				<li>Snow: {Snow}</li>
				<li>Rain: {Rain}</li>
		</ul></div>
);

class WeatherTable extends React.Component {
	
	state = {
		loading: true,
		weather: null,
	}	


	async componentDidMount() {
			//const url = "http://localhost:8080/api/24hr";
			const url = "/api/24hr";


			const resp = await fetch(url);
			const data = await resp.json();
			console.log(data);
		
			console.log(data);
			this.setState({loading: false, weather: data});
	}

	render() {
		if(this.state.loading || !this.state.weather) {
				return <div>Loading current weather...</div>;
		}

	

		return (
				<div>

				<div>
						<h3>24 Hour Stats</h3>
						<Stats High={this.state.weather.High} Low={this.state.weather.Low} Snow={this.state.weather.Snow} Rain={this.state.weather.Rain} />
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
