ALTER TABLE users
    RENAME CONSTRAINT users_email_key TO uni_users_email;

ALTER TABLE events
    RENAME CONSTRAINT events_user_id_fkey TO fk_events_user;
