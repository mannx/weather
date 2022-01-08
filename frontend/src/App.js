import React, {Component} from "react";
import "./App.css";

import WeatherTable from "./components/WeatherTable/WeatherTable.jsx";


class App extends Component {
	render() {
		return (
			<div className="App">
				<WeatherTable />
			</div>
		);
	}
}

export default App;
