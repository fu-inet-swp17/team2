$(document).ready(function(){
	$.ajax({
		url: "http://localhost/WebDashboard/chartjs/db_connection_datapoints.php",
		method: "GET",
		success: function(data) {
			for (var sensorTuple in data) {
				for (var sensorName in data[sensorTuple]) {
					for (var unitTuple in data[sensorTuple][sensorName]) {
						for (var unitName in data[sensorTuple][sensorName][unitTuple]) {
							var ID = [];
							var Value = [];
							var StringValue = [];
							var BoolValue = [];
							var DataValue = [];
							var Sum = [];
							var Time = [];
							var Link = [];
							
							for (var datapoint in data[sensorTuple][sensorName][unitTuple][unitName]) {
								ID.push(parseInt(data[sensorTuple][sensorName][unitTuple][unitName][datapoint].ID));
								var curValue;
								if ((curValue = data[sensorTuple][sensorName][unitTuple][unitName][datapoint].Value) != null) {
									Value.push(parseFloat(curValue));
								}
								var curStringValue;
								if ((curStringValue = data[sensorTuple][sensorName][unitTuple][unitName][datapoint].StringValue) != null) {
									StringValue.push(curStringValue.split(",").map(parseFloat));
								}
								var curBoolValue;
								if ((curBoolValue = data[sensorTuple][sensorName][unitTuple][unitName][datapoint].BoolValue) != null) {
									BoolValue.push(parseInt(curBoolValue));
								}
								var curDataValue;
								if ((curDataValue = data[sensorTuple][sensorName][unitTuple][unitName][datapoint].DataValue) != null) {
									DataValue.push(curDataValue);
								}
								Sum.push(data[sensorTuple][sensorName][unitTuple][unitName][datapoint].Sum);
								Time.push(new Date(data[sensorTuple][sensorName][unitTuple][unitName][datapoint].Time));
								Link.push(data[sensorTuple][sensorName][unitTuple][unitName][datapoint].Link);
							}
							
							createChart(sensorName, unitName, Value, StringValue, BoolValue, DataValue, Time);
						}
					}
				}
			}
					
				
		},
		error: function(data) {
			console.log(data);
		}
	});
});

function createChart(Name, Unit, Value, StringValue, BoolValue, DataValue, Time) {
	var FinalData = [];
	if (Value.length != 0) {
		FinalData = [Value];
	} else if (StringValue.length != 0) {
		FinalData = _.zip.apply(_, StringValue);
	} else if (BoolValue.length != 0) {
		FinalData = [BoolValue];
	} else if (DataValue.length != 0) {
		FinalData = DataValue;
	}
	
	document.getElementById("sections").innerHTML += "<div id='section_" + Name + "_" + Unit + "' style='height:500px; margin:100px 200px 100px 200px'></div>";
	google.charts.load('current', {'packages':['corechart']});
	google.charts.setOnLoadCallback(drawChart);

	function drawChart() {
		var data = new google.visualization.DataTable();
		data.addColumn('datetime', 'Time');
		var dataRow = [Time];
		for (var i=0; i<FinalData.length; i++) {
			data.addColumn('number', Unit);
			dataRow.push(FinalData[i]);
		}
		data.addRows(_.zip.apply(_, dataRow));
	
		var options = {
			title: Name,
			legend: { position: 'bottom' }
		};

		var chart = new google.visualization.LineChart(document.getElementById("section_" + Name + "_" + Unit));
		
		chart.draw(data, options);
	}
}
