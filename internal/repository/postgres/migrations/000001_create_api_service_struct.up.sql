CREATE TYPE lesson_type AS ENUM (
  'lecture',   -- лекция
  'seminar',   -- семинар
  'practice',  -- практика
  'test',      -- зачет
  'exam',      -- экзамен
  'webinar'    -- вебинар
);

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
  group_id integer,
  university_id INTEGER,
  FOREIGN KEY (university_id) REFERENCES universities(id),
  FOREIGN KEY (group_id) REFERENCES groups(id)
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

CREATE TABLE groups (
  id INTEGER PRIMARY KEY,
  name TEXT
);

CREATE TABLE universities_groups (
  group_id INTEGER,
  university_id INTEGER,
  FOREIGN KEY (group_id) REFERENCES groups(id),
  FOREIGN KEY (university_id) REFERENCES universities(id)
);

CREATE TABLE schedule (
  id INTEGER PRIMARY KEY,
  group_id INTEGER NOT NULL,
  subject TEXT NOT NULL,
  teacher TEXT,
  room TEXT,
  start_time TEXT NOT NULL, -- формат "09:30"
  end_time TEXT NOT NULL,   -- формат "11:00"
  date DATE,                -- конкретная дата (например, если пара по расписанию)
  day_of_week INTEGER,      -- если пара повторяется по дням недели
  type lesson_type NOT NULL, -- тип занятия (ENUM)
  created_at TIMESTAMP DEFAULT NOW()
);
