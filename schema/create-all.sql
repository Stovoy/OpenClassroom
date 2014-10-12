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

-- Foreign Keys

ALTER TABLE chat_messages ADD FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE chat_messages ADD FOREIGN KEY (chat_id) REFERENCES chats (id);
ALTER TABLE activity ADD FOREIGN KEY (user_id) REFERENCES users (id);