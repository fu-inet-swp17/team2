$(document).ready(function(){
	$.ajax({
		url: "http://localhost/scrollspy/chartjs/db_connection_final.php",
		method: "GET",
		success: function(data) {
			console.log(data);

			var ID = [[]];
			var Unit = [[]];
			var Value = [[]];
			var StringValue = [[]];
			var BoolValue = [[]];
			var DataValue = [[]];
			var Sum = [[]];
			var Time = [[]];
			var Link = [[]];
			
			var j = 0;
			var name = data[0].Name;
			for(var i in data) {
				if (data[i].Name != name) {
					createCanvas(name, Value[j], StringValue[j], BoolValue[j], DataValue[j], Time[j]);
					name = data[i].Name;
					j++;
					
					ID.push([]);
					Unit.push([]);
					Value.push([]);
					StringValue.push([]);
					BoolValue.push([]);
					DataValue.push([]);
					Sum.push([]);
					Time.push([]);
					Link.push([]);
				}
				ID[j].push(data[i].ID);
				Unit[j].push(data[i].Unit);
				Value[j].push(data[i].Value);
				StringValue[j].push(data[i].StringValue);
				BoolValue[j].push(data[i].BoolValue);
				DataValue[j].push(data[i].DataValue);
				Sum[j].push(data[i].Sum);
				Time[j].push(data[i].Time);
				Link[j].push(data[i].Link);
			}
			createCanvas(name, Value[j], StringValue[j], BoolValue[j], DataValue[j], Time[j]);
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