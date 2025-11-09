CREATE TABLE universities (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE,
  api_url TEXT,
  schedule bool
);

CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  chat_id TEXT UNIQUE, -- id чата/пользователя в Max
  name TEXT,
  group TEXT,
  university_id INTEGER,
  FOREIGN KEY (university) REFERENCES universities(id)
);

CREATE TABLE projects (
  id INTEGER PRIMARY KEY,
  name TEXT,
  description TEXT,
  university INTEGER
);

CREATE TABLE universities_projects (
  university_id INTEGER,
  project_id INTEGER,
  FOREIGN KEY (university_id) REFERENCES universities(id),
  FOREIGN KEY (project_id) REFERENCES projects(id)
);
