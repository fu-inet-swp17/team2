$(document).ready(function(){
	$.ajax({
		url: "http://localhost/SWP/chartjs/db_connection_chart_1.php",
		method: "GET",
		success: function(data) {
			console.log(data);
			var value = [];
			var BoardId = [];
			var value2 = [];
			var value3 = [];
			//first chart
			for(var i in data) {
				BoardId.push(data[i].ID);
				value.push(data[i].BoolValue);
				value2.push(data[i].Value);
				value3.push(data[i].Sum);
				
			}

			var chartdata = {
				labels: BoardId,
				datasets : [
					{
						label: 'SensorType',
						backgroundColor: "rgb(0,0,255)",
						data: value
					}
				]
			};

			var ctx = $("#mycanvas");

			var barGraph = new Chart(ctx, {
				type: 'line',
				data: chartdata
			});
			//end of first
			
			//second chart
			var chartdata = {
				labels: BoardId,
				datasets : [
					{
						label: 'Temperatur',
						backgroundColor: "rgb(224,255,0)",
						data: value
					}
				]
			};

			var ctx = $("#Temperatur");

			var barGraph = new Chart(ctx, {
				type: 'line',
				data: chartdata
			});
			//end of second chart
			
			//third chart
			var chartdata = {
				labels: BoardId,
				datasets : [
					{
						label: 'Luftdruck',
						backgroundColor: "rgb(224,0,255)",
						data: value2
					}
				]
			};

			var ctx = $("#Luftdruck");

			var barGraph = new Chart(ctx, {
				type: 'line',
				data: chartdata
			});
			//end of third chart
			
			//Fourth chart
			var chartdata = {
				labels: BoardId,
				datasets : [
					{
						label: 'Luftfeuchtigkeit',
						backgroundColor: "rgb(0,255,255)",
						data: value3
					}
				]
			};

			var ctx = $("#Luftfeuchtigkeit");

			var barGraph = new Chart(ctx, {
				type: 'bar',
				data: chartdata
			});
			//end of Fourth chart
			
			
			
			
		},
		error: function(data) {
			console.log(data);
		}
	});
});