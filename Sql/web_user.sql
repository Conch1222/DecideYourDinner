create database web;
use web;
create table web_user(
     user_id integer AUTO_INCREMENT not null,
     last_name varchar(255),
     first_name varchar(255),
     user_name varchar(255)  not null,
     password_hash varchar(70) not null,
     create_time datetime,
     primary key(user_id)
);

create table web_query_record(
     query_id integer AUTO_INCREMENT not null,
     user_id integer not null,
     store_name varchar(255) not null,
     store_address varchar(255),
     store_rating float4 not null,
     store_map_link varchar(255),
     create_time datetime,
     primary key (query_id),
     foreign key (user_id) REFERENCES web_user(user_id)
);

drop table web_user;
drop table web_query_record;

select * from web_user;
select * from web_query_record;
insert into web_user values (1, 'admin', 'admin', 'admin', '8C6976E5B5410415BDE908BD4DEE15DFB167A9C873FC4BB8A81F6F2AB448A918', '2024-05-18 17:02:00')