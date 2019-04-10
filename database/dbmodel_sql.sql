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
  user VARCHAR(16) REFERENCES user(id),
  create_time DATETIME,
  title VARCHAR(255),
INDEX create_time (create_time)
);

CREATE TABLE problem_entry (
  id VARCHAR(16) PRIMARY KEY,
  title VARCHAR(255),
  problem VARCHAR(16) REFERENCES problem(id)
);

CREATE TABLE submission (
  id VARCHAR(16) PRIMARY KEY,
  problem_judger VARCHAR(255),
  user VARCHAR(16) REFERENCES user(id),
  data BLOB
);


