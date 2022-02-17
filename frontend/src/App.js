import React, {Component} from "react";
import "./App.css";

import WeatherTable from "./components/WeatherTable/WeatherTable.jsx";
import Header from "./components/Header/Header.jsx";


class App extends Component {
	render() {
		return (
			<div className="App">
				<Header />
				<WeatherTable />
			</div>
		);
	}
}

export default App;
