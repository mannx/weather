import React from "react";

class WeatherTableTest extends React.Component {
	
	state = {
		loading: true,
		weather: null,
	}	
	
	async componentDidMount() {
			const url = "http://localhost:8080/api/test";
			const resp = await fetch(url);
			const data = await resp.json();
			console.log(data);
		
			this.setState({loading: false, weather: data});
	}

	render() {
		if(this.state.loading) {
				return <div>Loading current weather...</div>;
		}

		return (
				<div>
				<ul>
						<li>Current Temp: {this.state.weather.Temp}</li>
						<li>Feels Like: {this.state.weather.FeelsLike}</li>
						<li>Wind: {this.state.weather.WindSpeed}/{this.state.weather.WindDir}&deg;</li>
						<li>Rain: {this.state.weather.Rain1h}</li>
						<li>Snow: {this.state.weather.Snow1h}</li>
				</ul>
				</div>
		);
	}
}

export default WeatherTableTest;
