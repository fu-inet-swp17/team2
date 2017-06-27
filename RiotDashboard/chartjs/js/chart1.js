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
				BoardId.push("BoardId " + data[i].BoardId);
				value.push(data[i].Value1);
				value2.push(data[i].Value2);
				value3.push(data[i].Value3);
				
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
						label: 'Luftdruk',
						backgroundColor: "rgb(224,0,255)",
						data: value2
					}
				]
			};

			var ctx = $("#Luftdruk");

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
						label: 'Luftfr\u00fcchtigkeit',
						backgroundColor: "rgb(0,255,255)",
						data: value3
					}
				]
			};

			var ctx = $("#hi");

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