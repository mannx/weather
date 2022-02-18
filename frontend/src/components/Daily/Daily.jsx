import React from "react";
import UrlGet from "../URL/URL.jsx";

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
			errMsg = this.state.errMsg === null ? "Unknown Error" : this.state.errMsg;

			return <h3>Error has occurred: {errMsg}</h3>;
		}

		return <h3>Data here</h3>;
	}

	loadData = async () => {
		const month = this.state.date.getMonth()+1;		//month is 0 based
		const day = this.state.date.getDay();
		const year = this.state.date.getFullYear();

		const url = UrlGet("Daily") + "?month="+month+"&day="+day+"&year="+year;
		const resp = await fetch(url);
		const data = await resp.json();

		this.setState({error: false, errMsg: null});	// clear any previous errors

		if(data.Error !== null && data.Message !== null) {
			this.setState({error:true,errMsg:data.Message});
		}else{
			this.setState({data: data, loading: false});
		}
	}
}
