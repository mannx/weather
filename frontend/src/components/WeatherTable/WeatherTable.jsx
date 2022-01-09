import React from "react";

const User = ({Temp, WindSpeed, Snow1h}) => (
		<div><ul>
			<li>Temp: {Temp}</li>
			<li>Wind: {WindSpeed}</li>
			<li>Snow: {Snow1h}</li>
		</ul></div>
);

class WeatherTable extends React.Component {
	
	state = {
		loading: true,
		weather: null,
	}	


	async componentDidMount() {
			const url = "http://localhost:8080/api/24hr";
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
				{this.state.weather.map((w) => (
					<User Temp={w.Temp} WindSpeed={w.WindSpeed} Snow1h={w.Snow1h} />
				))}
			</div>
		);
	}
}

export default WeatherTable;
