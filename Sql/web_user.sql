create database web;
use web;
create table web_user(
	user_id integer not null,
	last_name varchar(255) not null,
    first_name varchar(255) not null,
    user_name varchar(255)  not null,
    password_hash varchar(70) not null,
    create_time datetime,
    primary key(user_id)
);

select * from web_user;
insert into web_user values (1, 'admin', 'admin', 'admin', '8C6976E5B5410415BDE908BD4DEE15DFB167A9C873FC4BB8A81F6F2AB448A918', '2024-05-18 17:02:00')