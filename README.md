Please add the following table on your local Postgre before running the webapp.

create table employees (
  id serial primary key,
  name varchar(256),
  team varchar(256)
);
