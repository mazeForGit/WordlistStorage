
class App extends React.Component {
	constructor(props) {
		super(props);
		this.state = { 
			pagetoreport: "impart.de",
			reportAsList: true,
			words: [{"id": 0, "name": "no word present", "occurance": 0, "new": false, "tests": null},],
			tests: {"test": "empty", "Category": [ { "name": "empty", "referencewordcount": 0, "referencecountsum": 0, "referenceoccurancesum": 0,"globalwordcount": 0,"globalcountsum": 0,"globaloccurancesum": 0,"localwordcount": 0,"localoccurancesum": 0}]}
		};
		this.handleChange = this.handleChange.bind(this);
		this.handleRun = this.handleRun.bind(this);
		this.handleChangeOnReportList = this.handleChangeOnReportList.bind(this);
	}
	handleChangeOnReportList(evt) {
		this.setState({ reportAsList: evt.target.checked });
	}
	handleChange(event) {
		this.setState({pagetoreport: event.target.value});
	}
	handleRun(event) {
		this.readData();
	}
	async readData() {
		try {
			console.log('WordList readData')
			var reqUrl = ""
			if (this.state.reportAsList) {
				if (window.location.port == "") {
					reqUrl = window.location.protocol + "//" + window.location.hostname + "/result?domain=" + this.state.pagetoreport;
				} else {
					reqUrl = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/result?domain=" + this.state.pagetoreport;
				}		
				const res = await fetch(reqUrl);
				const blocks = await res.json();
				
				if (blocks != null) {
					this.setState({
						words: blocks,
					})
				} else {
					var a = new Array()
					a = [{"id": 0, "name": "no mapping words found", "count": 0, "occurance": 0, "new": false, "tests": null},]
					this.setState.words = a
				}
			} else {
				if (window.location.port == "") {
					reqUrl = window.location.protocol + "//" + window.location.hostname + "/result?test=Big Five&domain=" + this.state.pagetoreport;
				} else {
					reqUrl = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port + "/result?test=Big Five&domain=" + this.state.pagetoreport;
				}		
				const res = await fetch(reqUrl);
				const blocks = await res.json();
				
				if (blocks != null) {
					this.setState({
						tests: blocks,
					})
				} else {
					var a = new Array()
					a = {"test": "empty", "Category": [ { "name": "empty", "referencewordcount": 0, "referencecountsum": 0, "referenceoccurancesum": 0,"globalwordcount": 0,"globalcountsum": 0,"globaloccurancesum": 0,"localwordcount": 0,"localoccurancesum": 0}]}
					this.setState.words = a
				}
			}
				
		} catch (e) {
			console.log(e);
		}
	}
	render() {
		return(
			<div className="container-fluid">
				<div className="container">
					<p></p>
					report on domain = &nbsp;
					<input type="text" size="40" value={this.state.pagetoreport} onChange={this.handleChange} />
					&nbsp;&nbsp;
					
					{(this.state.executionstarted && !this.state.executionfinished) ? 
						<button class="btn btn-primary  btn-sm" type="button" enabled onClick={this.handleRun}>
						  <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Processing
						</button>
					: 	<button class="btn btn-primary  btn-sm" type="button" enabled onClick={this.handleRun}>
							get report
						</button>
					}

					&nbsp;&nbsp; 
					<input type="checkbox" checked={this.state.reportAsList} onChange={this.handleChangeOnReportList} /> 
					&nbsp;as List
				</div>	
				
				{(this.state.reportAsList) ? <WordList words={this.state.words} /> : <TestList tests={this.state.tests} />}
		
			</div>
		);
	}
}
class TestList extends React.Component {
	constructor(props) {
		super(props);
	}
	render() {
		return (
			<div className="container">
				<p>Report on Test = Big Five</p>
				
				<div className="container">
					<div class="card-columns">
						{this.props.tests.Category.map(function(category, i) {
							return <Category key={i} category={category} />;
						})}
					</div>
				</div>
			</div>
		);
	}
}
class Category extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			words: []
		};	
	}
	render() {
		return (
			<div class="card">
				<div class="card-header">
					{this.props.category.name}
				</div>
				<div class="card-body">
					local: count={this.props.category.localwordcount} / occurance={this.props.category.localoccurancesum}
					<br />global: count={this.props.category.globalwordcount} / c-sum={this.props.category.globalcountsum} / o-sum={this.props.category.globaloccurancesum}
					<br />reference: count={this.props.category.referencewordcount} / c-sum={this.props.category.referencecountsum}
			
				</div>
			</div>
		)
	}
}
class WordList extends React.Component {
	constructor(props) {
		super(props);
	}
	render() {
		return (
			<div className="container">
				<p>list of words</p>
				
				<div className="container">
					<div class="card-columns">
						{this.props.words.map(function(word, i) {
							return <Word key={i} word={word} />;
						})}
					</div>
				</div>
			</div>
		);
	}
}
class Word extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			words: []
		};	
	}
	render() {
		return (
			<div class="card">
				<div class="card-header">
					count={this.props.word.occurance} (average={this.props.word.count})
				</div>
				<div class="card-body">
					{this.props.word.name}
				</div>
			</div>
		)
	}
}
ReactDOM.render(<App />, document.getElementById('app'));
