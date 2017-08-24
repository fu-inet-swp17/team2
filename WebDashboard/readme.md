# Prerequisites

## PHP server

Just set up a PHP server of your liking, e.g. Apache with XAMPP, and deploy all contents of this folder there. localhost should point to the contents of this folder. The website can be accessed through index.html.

## Path to PHP database connection file

Edit chartjs/js/chart.js and insert the domain of WebDashboard into the *domain* variable at the beginning.

## SQL database authentication

Edit the chartjs/config.php file to enter the authentication information for your SQL database. The database should be the one used by SenPi to store sensor data. If that database isn't available to you and you'd like to test the WebDashboard anyway, keep in mind your test database should conform to the correct database scheme (see wiki).

