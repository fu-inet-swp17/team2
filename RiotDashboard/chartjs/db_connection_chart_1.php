<?php


//setting header to json
header('Content-Type: application/json');

//config.php database and get connection  
include '/config.php';

//query to get data from the table
$query = sprintf("SELECT Value1,Value2,Value3,BoardId FROM DataTableName ");

//execute query
$result = $mysqli->query($query);

//loop through the returned data
$data = array();
foreach ($result as $row) {
	$data[] = $row;
}

//free memory associated with result
$result->close();

//close connection
$mysqli->close();

//now print the data
print json_encode($data);

?>