CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  name text not null,
  passport text default '',
  surname text not null,
  address text not null,
  patronomic text default ''
);
