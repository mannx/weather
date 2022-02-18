import React from "react";
import WeatherTable from "../WeatherTable/WeatherTable.jsx";
import Weekly from "../Weekly/Weekly.jsx";
import Daily from "../Daily/Daily.jsx";

import "./header.css";

function NavButton(props) {
	return <span className={"navLink"} onClick={props.func}>{props.name}</span>;
}

const Page24Hours = 0;
const PageWeekly = 1;
const PageDaily = 2;

function Navigate(props) {
	switch(props.page) {
		case Page24Hours:
			return <WeatherTable />;
		case PageWeekly:
			return <Weekly />;
		case PageDaily:
			return <Daily />;
		default:
			return <h1>Invalid Page Selected</h1>;
	}
}

export default class Header extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			page: Page24Hours,
		}
	}
		
	render() {
		return (<>
			<div><ul className={"navControl"}>
				<li className={"navControl"}><NavButton name={"24 Hours"} func={this.Nav24Hours} /></li>
				<li className={"navControl"}><NavButton name={"Weekly"} func={this.NavWeekly} /></li>
				<li className={"navControl"}><NavButton name={"Daily"} func={this.NavDaily} /></li>
			</ul></div>
			<Navigate page={this.state.page} />
		</>);
	}

	Nav24Hours = () => {this.setState({page: Page24Hours})}
	NavWeekly = () => {this.setState({page: PageWeekly})}
	NavDaily = () => {this.setState({page: PageDaily})}
}
