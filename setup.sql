drop table if exists photos;

create table photos (
  id int primary key,
  sol int,
  rover varchar(255),
  camera varchar(255),
  earthdate varchar(255),
  nasaimgsrc varchar(255),
  s3imgsrc varchar(255)
);

alter table photos owner to rover;
