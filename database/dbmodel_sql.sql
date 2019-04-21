CREATE TABLE user (
  id VARCHAR(16) PRIMARY KEY,
  user_name VARCHAR(255) UNIQUE,
  auth BLOB,
  register_time DATETIME,
INDEX register_time (register_time)
);

CREATE TABLE device (
  id VARCHAR(16) PRIMARY KEY,
  user VARCHAR(16) REFERENCES user(id),
  info BLOB
);

CREATE TABLE problem (
  id VARCHAR(16) PRIMARY KEY,
  problemset VARCHAR(16) REFERENCES problemset(id),
  user VARCHAR(16) REFERENCES user(id),
  create_time DATETIME,
  problem BLOB,
INDEX create_time (create_time)
);

CREATE TABLE problemset (
  id VARCHAR(16) PRIMARY KEY,
  user VARCHAR(16) REFERENCES user(id),
  problemset BLOB
);

CREATE TABLE submission (
  id VARCHAR(16) PRIMARY KEY,
  problem_judger VARCHAR(255),
  user VARCHAR(16) REFERENCES user(id),
  submission BLOB
);

CREATE TABLE judger (
  id VARCHAR(16) PRIMARY KEY,
  token VARCHAR(255)
);


