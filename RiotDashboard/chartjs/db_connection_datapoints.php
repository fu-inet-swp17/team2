<?php
//setting header to json
header('Content-Type: application/json');

//config.php database and get connection  
include '/config.php';

//get connection
$mysqli = new mysqli(DB_HOST, DB_USERNAME, DB_PASSWORD, DB_NAME);

if(!$mysqli){
	die("Connection failed: " . $mysqli->error);
}

$nameQuery = sprintf("SELECT DISTINCT Name FROM DataTable");
$nameResult = $mysqli->query($nameQuery);
$resultArray = array();
if (is_array($nameResult) || is_object($nameResult))
{
	// iterate over all different sensor names
	foreach ($nameResult as $name) {
		// get datapoints from that sepcific sensor
		$query = "SELECT * FROM DataTable WHERE Name = '".$name["Name"]."' ORDER BY Time";
		$result = $mysqli->query($query);
		
		// go over results from that specific datapoint
		$dataArray = array();
		if (is_array($result) || is_object($result)) {
			while($row = $result->fetch_assoc()) {
				array_push($dataArray, $row);
			}
			$result->close();
		}
		
		// add datapoints of the current sensor to the result array
		array_push($resultArray, array($name["Name"] => $dataArray));
	}
	
	$nameResult->close();
}

print json_encode($resultArray);

$mysqli->close();

?>