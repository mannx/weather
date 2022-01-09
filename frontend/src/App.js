import React, {Component} from "react";
import "./App.css";

import WeatherTable from "./components/WeatherTable/WeatherTable.jsx";


class App extends Component {
	render() {
		return (
			<div className="App">
				Weather Previous 24 hours:<br/>
				<WeatherTable />
			</div>
		);
	}
}

export default App;
