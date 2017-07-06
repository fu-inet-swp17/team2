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

//query to get data from the table
$query = sprintf("SELECT * FROM DataTable ORDER BY Name, Time");

//execute query
$result = $mysqli->query($query);

//loop through the returned data
$data = array();

if (is_array($result) || is_object($result))
{
	foreach ($result as $row) {
		array_push($data, $row);
	}
}

print json_encode($data);

//free memory associated with result
$result->close();

//close connection
$mysqli->close();

?>