# Cleanup
drop database insure;

# Create database
create database insure;

# change database
use insure;

create table binaries (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(200) NOT NULL, path VARCHAR(200) NOT NULL, email VARCHAR(200) NOT NULL, time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, is_virus_scanned ENUM('Y','N') NOT NULL DEFAULT 'N', vt_resourceid VARCHAR(200), vt_permalink VARCHAR(500), is_safe ENUM('Y','N'), PRIMARY KEY (id));

create table diversified_binaries (id INT NOT NULL AUTO_INCREMENT, binary_id INT NOT NULL, name VARCHAR(200) NOT NULL, path VARCHAR(200) NOT NULL, is_virus_scanned ENUM('Y','N') NOT NULL DEFAULT 'N', PRIMARY KEY (id), FOREIGN KEY (binary_id) REFERENCES binaries(id));

