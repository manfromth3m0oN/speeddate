CREATE TABLE users (
  id serial NOT NULL,
  PRIMARY KEY (id),
  email text NOT NULL,
  password text NOT NULL,
  name text NOT NULL,
  gender text NOT NULL,
  age integer NOT NULL,
  latitude real NOT NULL,
  longitude real NOT NULL,
  insert_ts timestamp default current_timestamp
);

CREATE TABLE swipe (
  id serial NOT NULL,
  PRIMARY KEY (id),
  swiper_id integer REFERENCES users (id) NOT NULL,
  swipee_id integer REFERENCES users (id) NOT NULL,
  preference boolean NOT NULL,
  insert_ts timestamp default current_timestamp
);


CREATE TABLE match (
  id serial NOT NULL,
  user_a integer REFERENCES users (id) NOT NULL,
  user_b integer REFERENCES users (id) NOT NULL,
  insert_ts timestamp default current_timestamp
);
