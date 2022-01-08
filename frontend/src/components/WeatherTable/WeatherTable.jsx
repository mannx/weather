import React from "react";

class WeatherTable extends React.Component {
	
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
				<div>Current Temp: {this.state.weather.Temp}</div>
		);
	}
}

export default WeatherTable;
