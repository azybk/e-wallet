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