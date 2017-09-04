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

// Select 25 players with lowest reaction time
$gameQuery = "SELECT * FROM ".DB_TABLE." WHERE NAME LIKE 'reaction_game%' ORDER BY Value ASC LIMIT 25";
$gameResult = $mysqli->query($gameQuery);
$resultArray = array();

if (is_array($gameResult) || is_object($gameResult)) {
  while($row = $gameResult->fetch_assoc()) {
    array_push($dataArray, $row);
  }
  $gameResult->close();
}

print json_encode($resultArray);

$mysqli->close();

?>
