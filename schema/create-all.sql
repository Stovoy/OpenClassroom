-- Table
CREATE TABLE users (
  id BIGSERIAL UNIQUE,
  name TEXT NOT NULL,
  lowername TEXT NOT NULL,
  encrypted_password TEXT NOT NULL,
  token TEXT NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE wiki (
  id BIGSERIAL UNIQUE,
  page TEXT NOT NULL,
  content TEXT NOT NULL
);

CREATE TABLE chats (
  id BIGSERIAL UNIQUE,
  wiki_id BIGINT NOT NULL
);

CREATE TABLE chat_messages (
  id BIGSERIAL UNIQUE,
  chat_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  time TIMESTAMP NOT NULL,
  message TEXT NOT NULL
);

CREATE TABLE activity (
  id BIGSERIAL UNIQUE,
  user_id BIGINT NOT NULL,
  action TEXT NOT NULL,
  time TIMESTAMP NOT NULL
);

CREATE TABLE lessons (
  id BIGSERIAL UNIQUE,
  name TEXT NOT NULL,
  teacher_id BIGINT NOT NULL,
  running BOOL NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP
);

CREATE TABLE lesson_items (
  id BIGSERIAL UNIQUE,
  lesson_id BIGINT NOT NULL,
  number BIGINT NOT NULL,
  action TEXT NOT NULL
);

-- Foreign Keys

ALTER TABLE chat_messages ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE chat_messages ADD FOREIGN KEY (chat_id) REFERENCES chats (id);
ALTER TABLE activity ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE lessons ADD FOREIGN KEY (teacher_id) REFERENCES users (id);
ALTER TABLE lesson_items ADD FOREIGN KEY (lesson_id) REFERENCES lessons (id);