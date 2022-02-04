import React from "react";
import {LineChart, Line, CartesianGrid, XAxis, YAxis, Tooltip} from "recharts";

class WeatherChart extends React.Component
{
	render() {
		console.log("wc: key: "+this.props.item);

		return (
			<LineChart width={600} height={300} data={this.props.data}>
					<Line type="monotone" dataKey={this.props.item} stroke="#8884d8" />
					<CartesianGrid stroke="#ccc" />
					<XAxis />
					<YAxis />
					<Tooltip />
			</LineChart>
		);
	}
}

export default WeatherChart;
