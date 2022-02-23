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
		}
	}

	componentDidMount() {
		this.loadData();
	}

	render() {
		if(this.state.loading === true || this.state.data === null) {
			return <h2>Loading data for day {this.state.date.toDateString()}</h2>;
		}

		if(this.state.error === true) {
			const errMsg = this.state.errMsg === null ? "Unknown Error" : this.state.errMsg;

			return <h3>Error has occurred: {errMsg}</h3>;
		}

		let wd = [];

		for(let k of Object.keys(this.state.data)) {
			wd.push([k, this.state.data[k]]);
		}

		return (<>
			{this.header()}
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
		const month = this.state.date.getMonth()+1;		//month is 0 based
		const day = this.state.date.getDate();
		const year = this.state.date.getFullYear();

		this.loadData2(month, day, year);
	}

	loadData2 = async (month, day, year) => {
		const url = UrlGet("Daily") + "?month="+month+"&day="+day+"&year="+year;
		const resp = await fetch(url);
		const data = await resp.json();

		this.setState({error: false, errMsg: null});	// clear any previous errors

		if(data.Error !== undefined && data.Message !== undefined) {
			this.setState({error:true,errMsg:data.Message});
		}else{
			this.setState({data: data, loading: false});
		}
	}


	header = () => {
		return (
			<div>
				<span>Pick Day to view stats:</span>
				<DatePicker selected={this.state.date} onChange={(e) => this.dateUpdated(e)} />
				{this.state.error && <span>{this.state.errMsg}</span>}
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
}
