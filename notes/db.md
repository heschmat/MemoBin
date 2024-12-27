# Setting up MySQL
```sh
sudo apt install mysql-server # ubuntu
```

## Scaffolding the database
```sh
sudo mysql

# Create a new UTF-8 `memobin` database:
CREATE DATABASE memobin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

SHOW DATABASES;

# Switch to using the `memobin` database:
USE memobin;

# Create a `memos` table.
CREATE TABLE memos (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

# Add an index on the create column.
CREATE INDEX idx_memos_created ON memos(created);

SHOW TABLES; # After selecting a database, list all tables in it.
DESCRIBE memos; # To check the structure of a specific table
DROP TABLE memos; # Drop the table.

# Add dummy records to the table.
INSERT INTO memos (title, content, created, expires) VALUES (
    'Python',
    'Simple, versatile, popular for data science and automation.',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 90 DAY)
);

INSERT INTO memos (title, content, created, expires) VALUES (
    'JavaScript',
    'Essential for dynamic websites and interactive web apps.',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 90 DAY)
);

INSERT INTO memos (title, content, created, expires) VALUES (
    'Go (Golang)',
    'Fast, scalable, great for modern backend systems.',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 90 DAY)
);

# Check user:
SELECT User, Host FROM mysql.user;
SELECT User, Host FROM mysql.user WHERE User = 'web' AND Host = 'localhost';

SELECT USER(), CURRENT_USER();

# Create a new `web` user (for security reason)
CREATE USER 'web'@'localhost';
DROP USER 'web'@'localhost';

CREATE USER 'web'@'localhost' IDENTIFIED BY 'changeme';
CREATE OR REPLACE USER 'web'@'localhost' IDENTIFIED BY 'changeme';

GRANT SELECT, INSERT, UPDATE, DELETE ON memobin.* TO 'web'@'localhost';
FLUSH PRIVILEGES; # apply the changes
# Set password (here: changeme):
ALTER USER 'web'@'localhost' IDENTIFIED BY 'changeme';



EXIT; # exit MySQL
mysql -u web -p # login as the `web` user
mysql -D memobin -u web -p # connect to `memobin` as the `web` user
SHOW GRANTS FOR CURRENT_USER();

INSERT INTO memos (title, content, created, expires) VALUES (
    'Ruby',
    'Elegant, readable, loved for rapid web development.',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 90 DAY)
);


# Deleting records --------------------------- #
DELETE FROM memos WHERE id BETWEEN 7 AND 10;

```

## Installing a database driver
To use MySQL from our Go web application we need to install a **database driver**. This essentially acts as a middleman, translating commands between Go and the MySQL database itself.

For our application we'll be using the popular `go-sql-driver/mysql` driver.

```sh
pwd # this should be `memobin` (your root project directory)

# get the latest version available, with the major release number 1
go get github.com/go-sql-driver/mysql@v1

```

## Setting up the session manager
```sh
sudo mysql;

USE memobin;
show tables;

SELECT user; # should be `root@localhost`

# `web@localhost` does NOT have the privilege to create table;
CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

```
token: a unique, randomly-generated, identifier for each session.

data: the actual session data we want to share between HTTP requests. (BLOB: binary large object)

expiry: expiry time for the session

# go mod
```sh
# To download the exact versions of all the packages that your project needs.
go mod download

# To ensure that nothing in the downloaded packages has been changed unexpectedly.
go mod verify

# To upgrade to latest available `minor` or `patch` release of a package
go get -u github.com/foo/bar

# To upgrade to a specific version
go get -u github.com/foo/bar@v3.2.1

# Removing unused packages -------------------
go get github.com/foo/bar@none

# or:
go mod tidy # automatically removes any unused packages from `go.mod` and `go.sum` files.
```

