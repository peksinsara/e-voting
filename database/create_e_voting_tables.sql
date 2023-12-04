-- Create the e-voting database
CREATE DATABASE IF NOT EXISTS e_voting;

-- Use the e-voting database
USE e_voting;

-- Create the User table
CREATE TABLE IF NOT EXISTS User (
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    full_name VARCHAR(255) NOT NULL,
    mothers_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    jmbg VARCHAR(13) NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT 0
);

-- Create the Candidate table
CREATE TABLE IF NOT EXISTS Candidate (
    candidate_id INT PRIMARY KEY AUTO_INCREMENT,
    full_name VARCHAR(255) NOT NULL,
    district VARCHAR(255) NOT NULL,
    short_bio VARCHAR(255) NOT NULL,
    total_votes INT DEFAULT 0
);

-- Create the Vote table
CREATE TABLE IF NOT EXISTS Vote (
    vote_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    candidate_id INT,
    FOREIGN KEY (user_id) REFERENCES User (user_id),
    FOREIGN KEY (candidate_id) REFERENCES Candidate (candidate_id)
);
