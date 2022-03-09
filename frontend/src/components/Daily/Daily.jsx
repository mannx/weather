import React from "react";
import NumberFormat from "react-number-format";
import DatePicker from "react-date-picker";
import UrlGet from "../URL/URL.jsx";

import "./style.css";

export default class Daily extends React.Component {
	constructor(props) {
		super(props);

		this.state = {
			date: new Date(),
			data: null,
			loading: true,
			error: false,
			errMsg: null,

			cities: [],		//list of cities we might have data for {ID: xxxx, Name: city_name}
			cityid: null,	// id of the city to display data for
		}
	}

	componentDidMount() {
		this.loadData();
	}

	render() {
		let output = null;

		if(this.state.cityid === null) {
			// no city id, display a message to select a city
			output = <></>;
		}else{
			// have a city id, return the data for it
			output = this.renderData();
		}

		return (<>
			{this.header()}
			{output}
		</>);
	}

	renderData = () => {
		if(this.state.data === null) {
			return <h3>Loading data</h3>;
		}

		let wd = [];

		for(let k of Object.keys(this.state.data)) {
			wd.push([k, this.state.data[k]]);
		}

		return (<>
			<table>
				<thead><tr>
					<th>Item</th>
					<th>Value</th>
				</tr></thead>
				<tbody>
					{wd.map(function(obj, i){
						return (
							<tr key={i}><td>{obj[0]}</td>
								<td><NumberFormat decimalScale={2} value={obj[1]} displayType="text" thousandSeparator="true" fixedDecimalScale="true" /></td>
							</tr>
						);
					})}
				</tbody>
			</table>
		</>);
	}

	loadData = async () => {
		// get the list of cities if we dont have them
		if(this.state.cities === undefined || this.state.cities === null || this.state.cities.length === 0) {
			const url = UrlGet("CityList");
			const resp = await fetch(url);
			const data = await resp.json();

			this.setState({cities: data});
		}

		const month = this.state.date.getMonth()+1;		//month is 0 based
		const day = this.state.date.getDate();
		const year = this.state.date.getFullYear();

		this.loadData2(month, day, year, 0);
	}

	loadData2 = async (month, day, year, city) => {
		const cid = city === undefined ? 0 : city;
		const url = UrlGet("Daily") + "?month="+month+"&day="+day+"&year="+year+"&city="+cid;
		const resp = await fetch(url);
		const data = await resp.json();

		this.setState({error: false, errMsg: null});	// clear any previous errors

		if(data.Error !== undefined && data.Message !== undefined) {
			this.setState({error:true,errMsg:data.Message});
			this.setState({data:null, loading:true});
		}else{
			this.setState({data: data, loading: false});
		}
	}

	loadDaily = async (city) => {
		const month = this.state.date.getMonth()+1;		//month is 0 based
		const day = this.state.date.getDate();
		const year = this.state.date.getFullYear();
		const cid = city === undefined ? 0 : city;

		this.loadData2(month, day, year, cid);
	}

	header = () => {
		return (
			<div>
				<div>
				<span>Pick Day to view stats:</span>
				<DatePicker selected={this.state.date} onChange={(e) => this.dateUpdated(e)} />
				</div>
				<div>
				<label>Pick city to view
				<select value={this.state.selectedCity} onChange={(e) => this.cityUpdated(e)}>
					<option value={0}>Pick a city</option>
					{this.state.cities.map(function(obj, i) {
						return <option value={obj.ID}>{obj.Name}</option>;
					})}
				</select>
				</label>
				</div>
				{this.state.error && <div>{this.state.errMsg}</div>}
			</div>
		);
	}

	dateUpdated = (e) => {
		this.setState({date: e});
		this.setState({error: false, errMsg: null});

		// state doesnt update immediatly
		const month = e.getMonth() + 1;
		const year = e.getFullYear();
		const day = e.getDate();

		this.loadData2(month, day, year);
	}

	cityUpdated = (e) => {
		this.setState({cityid: e.target.value});
		this.loadDaily(e.target.value);
	}
}
