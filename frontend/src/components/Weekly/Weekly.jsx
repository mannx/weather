import React from "react";
import DatePicker from "react-date-picker";
import UrlGet from "../URL/URL.jsx";


// used to display data for the previous week, or a selected week
export default class Weekly extends React.Component {
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
		const title = <h3>Weekly data for week ending: {this.state.date.toDateString()}</h3>;
		const emsg = this.state.error ? this.state.errMsg : null;
		const error = <span className="error" >{emsg}</span>;

		return (<>
			{title}
			{error}
		</>);
	}

	loadData = async () => {
		const year = this.state.date.getFullYear();
		const month = this.state.date.getMonth() + 1;		// 0 based return
		const day = this.state.date.getDay();

		const url = UrlGet("Weekly") + "?month="+month+"&year="+year+"&day="+day;
		const resp = await fetch(url);
		const data = await resp.json();

		let err = false;
		let errmsg = null;

		if(data.Error !== null) {
			err = true;
			errmsg = data.Message;
		}else{
			this.setState({data: data, loading: false});
		}

		this.setState({error:err,errMsg: errmsg});
	}
}
