dependecies

go get github.com/gofiber/fiber/v2
go get -u github.com/doug-martin/goqu/v9
go get github.com/lib/pq
go get golang.org/x/crypto/bcrypt
go get github.com/allegro/bigcache
go get github.com/joho/godotenv

create database e_wallet;
create table users
(
id serial primary key,
full_name varchar(255),
phone varchar(20),
username varchar(30),
password varchar(30)
)

ALTER TABLE users
ADD COLUMN email_verified_at TIMESTAMP
ADD COLUMN email varchar(100)

create table account
(
id serial primary key,
user_id int,
account_number varchar(50),
balance real
)

create table transaction
(
id serial primary key,
account_id int,
sof_number varchar(50),
dof_number varchar(50),
transaction_type char,
amount real,
transaction_datetime timestamp
)

go get github.com/redis/go-redis/v9

CREATE TABLE notifications
(
id serial PRIMARY KEY,
user_id INT NOT NULL,
title TEXT NOT NULL,
body TEXT NOT NULL,
status INT NOT NULL,
is_read INT NOT NULL,
created_at TIMESTAMP NOT NULL
)

create table template
(
code varchar(100) PRIMARY KEY,
title varchar(100),
body text
)

create table topup
(
id varchar(100) not null primary key,
user_id int,
amount real,
status int default 0,
snap_url varchar(255)
)
