CREATE TABLE users_tasks(
  id SERIAL PRIMARY KEY,
  user_id int REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  start_date timestamp not null,
  end_date timestamp not null,
  last_start timestamp not null,
  duration int not null
);
