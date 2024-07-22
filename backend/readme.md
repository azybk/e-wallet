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
