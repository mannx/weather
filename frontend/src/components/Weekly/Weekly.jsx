import React from "react";


// used to display data for the previous week, or a selected week
export default class Weekly extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			date: new Date(),
		}
	}

	render() {
		return (
			<h3>Weekly data here, ending {date}</h3>;
		);
	}
}
