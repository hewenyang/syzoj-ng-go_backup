CREATE TABLE user (
  id VARCHAR(16) PRIMARY KEY,
  user_name VARCHAR(64) UNIQUE,
  auth BLOB
);

CREATE TABLE device (
  id VARCHAR(16) PRIMARY KEY,
  user VARCHAR(16) REFERENCES user(id),
  info BLOB
);

CREATE TABLE problem (
  id VARCHAR(16) PRIMARY KEY,
  title VARCHAR(255),
  user VARCHAR(16) REFERENCES user(id)
);

CREATE TABLE problem_entry (
  id VARCHAR(16) PRIMARY KEY,
  problem VARCHAR(16) REFERENCES problem(id)
);

CREATE TABLE problem_source (
  id VARCHAR(16) PRIMARY KEY,
  source BLOB
);

CREATE TABLE problem_judger (
  id VARCHAR(16) PRIMARY KEY,
  problem VARCHAR(16) REFERENCES problem(id),
  user VARCHAR(16) REFERENCES user(id),
  type VARCHAR(255),
  data BLOB
);

CREATE TABLE problem_statement (
  id VARCHAR(16) PRIMARY KEY,
  problem VARCHAR(16) REFERENCES problem(id),
  user VARCHAR(16) REFERENCES user(id),
  data BLOB
);

CREATE TABLE submission (
  id VARCHAR(16) PRIMARY KEY,
  problem_judger VARCHAR(255),
  user VARCHAR(16) REFERENCES user(id),
  data BLOB
);


