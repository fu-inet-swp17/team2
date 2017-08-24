<?php
//setting header to json
header('Content-Type: application/json');

//config.php database and get connection  
include 'config.php';

//get connection
$mysqli = new mysqli(DB_HOST, DB_USERNAME, DB_PASSWORD, DB_NAME);
mysqli_set_charset($mysqli, "utf8");

if(!$mysqli) {
	die("Connection failed: " . $mysqli->error);
}

$nameQuery = "SELECT DISTINCT Name FROM ".DB_TABLE;
$nameResult = $mysqli->query($nameQuery);
$resultArray = array();
if (is_array($nameResult) || is_object($nameResult)) {
	// iterate over all different sensor names
	foreach ($nameResult as $name) {
		// get all units for that sensor
		$unitQuery = "SELECT DISTINCT Unit FROM ".DB_TABLE." WHERE Name = '".$name["Name"]."' ORDER BY Unit ASC";
		$unitResult = $mysqli->query($unitQuery);
		$innerResultArray = array();
		if (is_array($unitResult) || is_object($unitResult)) {
			// iterate over all different units for that sensor
			foreach ($unitResult as $unit) {
				// get datapoints from that sepcific sensor
				$query = "SELECT * FROM ".DB_TABLE." WHERE Name = '".$name["Name"]."' AND Unit = '".$unit["Unit"]."' ORDER BY Time";
				$result = $mysqli->query($query);
				
				// go over results from that specific datapoint
				$dataArray = array();
				if (is_array($result) || is_object($result)) {
					while($row = $result->fetch_assoc()) {
						array_push($dataArray, $row);
					}
					$result->close();
				}
				
				// add datapoints of the current unit to the result array
				array_push($innerResultArray, array($unit["Unit"] => $dataArray));
			}
			
			$unitResult->close();
		}
		
		// add datapoints of the current sensor to the result array
		array_push($resultArray, array($name["Name"] => $innerResultArray));
	}
	
	$nameResult->close();
}

print json_encode($resultArray);

$mysqli->close();

?>