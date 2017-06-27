<?php

//database
define('DB_HOST', '');//database host 
define('DB_USERNAME', '');//database username
define('DB_PASSWORD', '');//database password
define('DB_NAME', '');//database name

//get connection
$mysqli = new mysqli(DB_HOST, DB_USERNAME, DB_PASSWORD, DB_NAME);

if(!$mysqli){
	die("Connection failed: " . $mysqli->error);
}

?>