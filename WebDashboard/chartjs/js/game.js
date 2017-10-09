$(document).ready(function(){
	var domain = ""; // INSERT DOMAIN OF WEBDASHBOARD HERE (WITHOUT HTTP://)
	$.ajax({
		url: "http://" + domain + "/chartjs/game_results.php",
		method: "GET",
		success: function(data) {
			for (var result in data) {
        result = "<tr><td>".concat(result.Name, "</td><td>", result.Value, "</td></tr>");
        $('#resultTable tr:last').after(result);
			}
		},
		error: function(data) {
			console.log(data);
		}
	});
});
