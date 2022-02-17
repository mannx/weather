import React from "react";

import "./header.css";

function NavButton(props) {
	return <span className={"navLink"} onClick={props.func}>{props.name}</span>;
}

const Page24Hours = 0;
const PageWeekly = 1;

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
			</ul></div>
		</>);
	}

	Nav24Hours = () => {}
	NavWeekly = () => {}
}
