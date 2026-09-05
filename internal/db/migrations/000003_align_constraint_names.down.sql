ALTER TABLE events
    RENAME CONSTRAINT fk_events_user TO events_user_id_fkey;

ALTER TABLE users
    RENAME CONSTRAINT uni_users_email TO users_email_key;
