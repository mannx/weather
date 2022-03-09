import React from "react";
import UrlGet from "../URL/URL.jsx";

export default class UpdateConfig extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			error: false,
			message: "Place Holder Message",
		}
	}

	render() {
		return(<>
			<fieldset>
				<legend>Migrate Configuration File</legend>
				{this.message()}
				<button onClick={this.migrate}>Migrate</button>
			</fieldset>
		</>);
	}

	migrate = async () => {
		const url = UrlGet("Migrate");
		const resp = await fetch(url);
		const data = await resp.json();

		this.setState({error: data.Error, message: data.Message})
	}

	message = () => {
		if(this.state.message === null){
			return <></>;
		}

		const colour = this.state.error ? "red" : "blue";

		return <span style={{"color": colour}}>{this.state.message}</span>;
	}
}
