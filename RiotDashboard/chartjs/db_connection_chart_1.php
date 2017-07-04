<?php
//setting header to json
header('Content-Type: application/json');

//database
define('DB_HOST', 'localhost');
define('DB_USERNAME', 'root');
define('DB_PASSWORD', 'SoftwareProject');
define('DB_NAME', 'test');

//get connection
$mysqli = new mysqli(DB_HOST, DB_USERNAME, DB_PASSWORD, DB_NAME);

if(!$mysqli){
	die("Connection failed: " . $mysqli->error);
}
/*
mysqli_query($mysqli,"CREATE TABLE DeviceTable (
					ID INT PRIMARY KEY AUTO_INCREMENT, 
					Address varchar(255) NOT NULL, 
					FailedAttempts INT DEFAULT 0);");
 

	
mysqli_query($mysqli,"CREATE TABLE  DataTable (ID INT PRIMARY KEY AUTO_INCREMENT,
					 Name VARCHAR(255) NULL DEFAULT NULL,
					 Unit VARCHAR(20) NULL DEFAULT NULL,
					 Value FLOAT NULL DEFAULT NULL,
					 StringValue VARCHAR(255) NULL DEFAULT NULL,
					 BoolValue BOOL NULL DEFAULT NULL,
					 DataValue VARCHAR(255) NULL DEFAULT NULL,
					 Sum FLOAT NULL DEFAULT NULL,
					 Time TIMESTAMP NULL DEFAULT NULL,
					 Link VARCHAR(255) NULL DEFAULT NULL);");


*/
//query to get data from the table
$nameQuery = sprintf("SELECT DISTINCT Name FROM DataTable ");
$nameResult = $mysqli->query($nameQuery);
if (is_array($nameResult) || is_object($nameResult))
{
	foreach ($nameResult as $nameRow) {
		$query = sprintf("SELECT * FROM DataTable WHERE Name='%s'", $nameRow["Name"]);
		
		//execute query
		$result = $mysqli->query($query);

		//loop through the returned data
		$data = array();
		if (is_array($result) || is_object($result))
		{
			foreach ($result as $row) {
				$data[] = $row;
			}
		}

		//now print the data
		print json_encode($data);
	}
}
//free memory associated with result
//$nameResult->close();

//close connection
$mysqli->close();
?>