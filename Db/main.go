package Db

var Schemas = []string{
	`CREATE TABLE IF NOT EXISTS Donator (id INT PRIMARY KEY AUTO_INCREMENT, name CHAR(50), bloodgroup VARCHAR(10), email VARCHAR(255), password VARCHAR(500), address VARCHAR(200), credits INT DEFAULT 0, phoneno VARCHAR(12), createdAt DATETIME DEFAULT CURRENT_TIMESTAMP());`,
	`CREATE TABLE IF NOT EXISTS Hospital (id INT PRIMARY KEY AUTO_INCREMENT, name CHAR(50), email VARCHAR(255), password VARCHAR(500), address VARCHAR(200), phoneno VARCHAR(12), createdAt DATETIME DEFAULT CURRENT_TIMESTAMP());`,
	`CREATE TABLE IF NOT EXISTS Admins (id INT PRIMARY KEY, name CHAR(50), createdAt DATETIME DEFAULT CURRENT_TIMESTAMP());`,
	`CREATE TABLE IF NOT EXISTS donations (id INT PRIMARY KEY, donorId INT REFERENCES donator(id), hospitalId INT REFERENCES hospital(id), createdAt DATETIME DEFAULT CURRENT_TIMESTAMP());`,
	`CREATE TABLE IF NOT EXISTS requests (id INT PRIMARY KEY, hospitalId INT REFERENCES hospital(id), type tinyint, bloodgroup VARCHAR(10));`,
}

type Hospital struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	Address   string `db:"address"`
	PhoneNo   string `db:"phoneno"`
	CreatedAt string `db:"createdAt"`
}

type Donator struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	BloodGroup string `db:"bloodgroup"`
	Email      string `db:"email"`
	Password   string `db:"password"`
	Address    string `db:"address"`
	PhoneNo    string `db:"phoneno"`
	CreatedAt  string `db:"createdAt"`
}
