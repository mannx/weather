import React from "react";
import {LineChart, Line, CartesianGrid, XAxis, YAxis, Tooltip} from "recharts";

class WeatherChart extends React.Component
{
	render() {
		return (
			<LineChart width={600} height={300} data={this.props.data}>
					<Line type="monotone" dataKey={this.props.item} stroke="#8884d8" />
					<CartesianGrid stroke="#ccc" />
					<XAxis dataKey="TimeString"/>
					<YAxis />
					<Tooltip />
			</LineChart>
		);
	}
}

export default WeatherChart;
