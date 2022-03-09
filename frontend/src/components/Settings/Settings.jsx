import React from "react";
import Popup from "reactjs-popup";
import UrlGet from "../URL/URL.jsx";

import UpdateConfig from "./Config.jsx";

import "reactjs-popup/dist/index.css";

export default class Settings extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			error: false,
			errMsg: null, 

			cities: [],			// list of cities we are monitoring
			newCityText: "",
			cityPopup: false,		// is the confirm city popup displayed?
			cityValidate: [],		// list of possible cities to add from user input
			validateCityID: null,	// id we are about to confirm
		}
	}

	async componentDidMount() {
		// load the cities we are monitoring
		this.loadCities();
	}

	render() {
		return (<>
			{this.header()}
			<UpdateConfig />
			{this.cityInfo()}
			{this.confirmCityPopup()}
		</>);
	}

	loadCities = async () => {
		const url = UrlGet("CityList");
		const resp = await fetch(url);
		const data = await resp.json();

		if(data !== null) {
			this.setState({cities: data});
		}
	}

	// display the settings header
	header = () => {
		return <h3>Settings</h3>;
	}

	// Display the current cities, add/remove cities we are gathering data for
	cityInfo = () => {
		return (<div>
			<fieldset>
				<legend>Cities</legend>
				{this.displayErrors()}
				<label>Add New City: <input type="text" value={this.state.newCityText} onChange={(e) => this.setState({newCityText: e.target.value})} /></label>
				<button onClick={this.addCity} >Add</button><button onClick={() => this.setState({newCityText: ""})}>Clear</button>
				<div>
					<ul className="citySettingList">
						{this.state.cities.map(function(obj, i) {
							//return <li>{obj.Name}</li>;
							return this.cityLink(obj);
						}, this)}
					</ul>
				</div>
			</fieldset>
		</div>);
	}

	// render the controls for display and removal of a city
	cityLink = (obj) => {
		return (<li className="citySettingList" key={obj.ID}>
			<input type="checkbox"/>{obj.Name}
		</li>);
	}

	displayErrors = () => {
		let err = null;

		if(this.state.error) {
			err = <span style={{color:"red"}}>Error: {this.state.errMsg}</span>;
		}

		return <div>{err}</div>;
	}

	// add a city to the database by name, possibly returns an error
	// this sends the city name, and returns validation data to confirm?
	// this could/should be done better, todo at some point
	addCity = () => {
		const url = UrlGet("CityAdd");
		const options = {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({city: this.state.newCityText}),
		}

		this.setState({cityPopup: true});

		fetch(url, options)
			.then(r => r.json())
			.then(r => this.setState({cityValidate: r, cityPopup: true}));
	}

	// display a popup with all the cities that were found to confirm which city to add, used if we have multiple returns
	confirmCityPopup = () => {
		return (<>
			<Popup open={this.state.cityPopup} onClose={() => {console.log("city confirmed")}}>
				<div className="modal"><div>
					<span>Select correct city:</span><br/>
					{this.state.cityValidate.map((obj) => {
						let state = "";
						if(obj.State || obj.State !== "") {
							state = ", " + obj.State;
						}

						return (<>
							<input key={obj.ID} type="radio" id={obj.ID} name="cityConfirm" value={obj.ID} onChange={(e) => this.setState({validateCityID: e.target.value})}/>
							<label >{obj.Name}{state}, {obj.Country}</label>
							<br/>
						</>);
					})}</div>
					<button onClick={this.confirmCity}>Confirm</button>
					<button onClick={()=>{this.setState({cityPopup: !this.state.cityPopup})}} >Close</button>
				</div>
			</Popup>
		</>);
	}

	confirmCity = () => {
		console.log("confirming city selection: " + this.state.validateCityID);
		this.setState({cityPopup: false});			// close the popup

		// add the city to the database and reload
		const url = UrlGet("CityConfirm");
		const options = {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({id: this.state.validateCityID})
		}

		fetch(url, options)
			.then(r => r.json())
			//.then(r => this.handleErrors(r))
			.then(r => {
				if(r.Error) {
					this.setState({error:true,errMsg: r.Message});
				}else{
					// update the city list
					this.setState({cities: r});
				}
			})
	}

	// sets the appropriate fields to display an error message if required, otherwise clears it
	handleErrors = (e) => {
		const msg = e.Error ? e.Message : "";
		this.setState({error: e.Error, errMsg: msg});
	}
}
