CREATE TABLE messages (
  id      BIGSERIAL PRIMARY KEY,
  body    text    NOT NULL,
  send_at text    NOT NULL,
  was_sent_at  text    NOT NULL,
  status  text    NOT NULL,
  tries   int NOT NULL
);
