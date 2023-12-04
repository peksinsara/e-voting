#!/bin/bash

# MySQL connection parameters
MYSQL_USER="root"
MYSQL_PASSWORD="2309"

# SQL script file
SQL_SCRIPT="create_e_voting_tables.sql"

# MySQL command to execute the script
MYSQL_COMMAND="mysql -u$MYSQL_USER -p$MYSQL_PASSWORD < $SQL_SCRIPT"

# Execute the SQL script
echo "Creating e-voting database and tables..."
$MYSQL_COMMAND

echo "Script execution completed."
