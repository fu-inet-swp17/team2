$(document).ready(function(){
	$.ajax({
		url: "http://localhost/scrollspy/chartjs/db_connection_datapoints.php",
		method: "GET",
		success: function(data) {
			for (var sensorTuple in data) {
				for(var sensorName in data[sensorTuple]) {
					var ID = [];
					var Unit = [];
					var Value = [];
					var StringValue = [];
					var BoolValue = [];
					var DataValue = [];
					var Sum = [];
					var Time = [];
					var Link = [];
					
					for(var datapoint in data[sensorTuple][sensorName]) {
						ID.push(data[sensorTuple][sensorName][datapoint].ID);
						Unit.push(data[sensorTuple][sensorName][datapoint].Unit);
						Value.push(data[sensorTuple][sensorName][datapoint].Value);
						StringValue.push(data[sensorTuple][sensorName][datapoint].StringValue);
						BoolValue.push(data[sensorTuple][sensorName][datapoint].BoolValue);
						DataValue.push(data[sensorTuple][sensorName][datapoint].DataValue);
						Sum.push(data[sensorTuple][sensorName][datapoint].Sum);
						Time.push(data[sensorTuple][sensorName][datapoint].Time);
						Link.push(data[sensorTuple][sensorName][datapoint].Link);
					}
					
					createCanvas(sensorName, Value, StringValue, BoolValue, DataValue, Time);
				}
			}
					
				
		},
		error: function(data) {
			console.log(data);
		}
	});
});

function createCanvas(name, Value, StringValue, BoolValue, DataValue, Time) {
	var FinalData = [];
	if (Value != []) {
		FinalData = Value;
	} else if (StringValue != []) {
		FinalData = StringValue;
	} else if (BoolValue != []) {
		FinalData = BoolValue;
	} else if (DataValue != []) {
		FinalData = DataValue;
	}
	
	document.getElementById("sections").innerHTML += "<div id='section_" + name + "'><div id='chart-container_" + name + "'><canvas id='" + name + "'" + "></canvas></div></div>";
	
	//create chart
	var chartdata = {
		labels: Time,
		datasets : [
			{
				label: name,
				backgroundColor: "rgb(0,0,255)",
				data: FinalData
			}
		]
	};

	var ctx = document.getElementById(name).getContext('2d');
	
	var lineGraph = new Chart(ctx, {
		type: 'line',
		data: chartdata
	});
	
	//end of chart
	
}