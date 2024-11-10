package Db

import "database/sql"

var Schemas = []string{
	`CREATE TABLE IF NOT EXISTS user (id INT PRIMARY KEY  AUTO_INCREMENT, email VARCHAR(255), password VARCHAR(500), type tinyint, createdAt DATETIME DEFAULT CURRENT_TIMESTAMP())`,
	`CREATE TABLE IF NOT EXISTS donator (userId INT REFERENCES user(id), name CHAR(50), bloodgroup VARCHAR(10), address VARCHAR(200), pincode VARCHAR(6), credits INT DEFAULT 0, phoneno VARCHAR(12));`,
	`CREATE TABLE IF NOT EXISTS hospital (userId INT REFERENCES user(id), name CHAR(50), address VARCHAR(200), pincode VARCHAR(6), phoneno VARCHAR(12));`,
	`CREATE TABLE IF NOT EXISTS admins (id INT PRIMARY KEY AUTO_INCREMENT, name CHAR(50), createdAt DATETIME DEFAULT CURRENT_TIMESTAMP());`,
	`CREATE TABLE IF NOT EXISTS donations (id INT PRIMARY KEY AUTO_INCREMENT, donorId INT REFERENCES donator(userId), hospitalId INT REFERENCES hospital(userId), createdAt DATETIME DEFAULT CURRENT_TIMESTAMP());`,
	`CREATE TABLE IF NOT EXISTS requests (id INT PRIMARY KEY AUTO_INCREMENT, hospitalId INT REFERENCES hospital(userId), type tinyint, bloodgroup VARCHAR(10), unit int, acceptedBy INT REFERENCES donator(userId));`,
	`CREATE TABLE IF NOT EXISTS blacklist (donatorId INT REFERENCES donator(userId), hospitalId INT REFERENCES Hospital(userId), createdAt DATETIME DEFAULT CURRENT_TIMESTAMP());`,
}

type Request struct {
	Id         int           `db:"id"`
	HospitalId int           `db:"hospitalId"`
	Type       int           `db:"type"`
	BloodGroup string        `db:"bloodgroup"`
	Unit       int           `db:"unit"`
	AcceptedBy sql.NullInt64 `db:"acceptedBy"`
}

type User struct {
	Id        int    `db:"id"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	Type      string `db:"type"`
	CreatedAt string `db:"createdAt"`
}

type Hospital struct {
	UserId    int    `db:"userId"`
	Name      string `db:"name"`
	Address   string `db:"address"`
	Pincode   string `db:"pincode"`
	PhoneNo   string `db:"phoneno"`
	CreatedAt string `db:"createdAt"`
}

type Donator struct {
	UserId     int    `db:"userId"`
	Name       string `db:"name"`
	BloodGroup string `db:"bloodgroup"`
	Address    string `db:"address"`
	Pincode    string `db:"pincode"`
	PhoneNo    string `db:"phoneno"`
	CreatedAt  string `db:"createdAt"`
}
