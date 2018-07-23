DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS follows;

CREATE TABLE public.users
(
    id serial NOT NULL,
    "name" text NULL,
    username text NOT NULL,
    email text NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    password_hash text null,
    remember_hash text null,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)
WITH (
	OIDS=FALSE
) ;
CREATE UNIQUE INDEX uix_users_email ON public.users USING btree
(email) ;
CREATE UNIQUE INDEX uix_users_username ON public.users USING btree
(username) ;

-- Insert Users
INSERT INTO public.users
    ("name", username, email, password_hash, remember_hash)
VALUES
    ('Sam Smith', 'samsmith', 'sam2018@gmail.com', 'fake-pw-hash', 'fake-hash');

INSERT INTO public.users
    ("name", username, email, password_hash, remember_hash)
VALUES
    ('Kanye West', 'kanye_west', 'kanye@kanye.com', 'fake-pw-hash', 'fake-hash');

INSERT INTO public.users
    ("name", username, email, password_hash, remember_hash)
VALUES
    ('Dua Lipa', 'duasings', 'dua@lipa.com', 'fake-pw-hash', 'fake-hash');

INSERT INTO public.users
    ("name", username, email, password_hash, remember_hash)
VALUES
    ('Bob Dylan', 'bobbyd', 'bob@dylan.com', 'fake-pw-hash', 'fake-hash');

-- this user's remember_hash based on
-- HMACKey: "secret-hmac-key" 
-- Used to test logout
INSERT INTO public.users
    ("name", username, email, password_hash, remember_hash)
VALUES
    ('Tom Tester', 'tommyTesterton', 'tommy@gmail.com', 'fake-pw-hash',
        'M-QQp9qW_mhAZK5s1651k_ODELJZTP1EBahOOByjfJE=');
-- this user's remember_hash based on
-- HMACKey: "secret-hmac-key" 
INSERT INTO public.users
    ("name", username, email, password_hash, remember_hash)
VALUES
    ('Main Tester', 'mainTester', 'mtester@gmail.com', 'fake-pw-hash',
        'Dt0b9x7U0tO22dNEX3f1uLMd5STOl5hbDU2ATW6pMjw=');

-- Drop table

-- DROP TABLE public.follows

CREATE TABLE public.follows
(
    follower_id serial NOT NULL,
    user_id serial NOT NULL,
    CONSTRAINT follows_pkey PRIMARY KEY (follower_id, user_id)
)
WITH (
	OIDS=FALSE
) ;

INSERT INTO public.follows
    (follower_id, user_id)
VALUES(1, 4);
INSERT INTO public.follows
    (follower_id, user_id)
VALUES(2, 4);
INSERT INTO public.follows
    (follower_id, user_id)
VALUES(3, 2);
INSERT INTO public.follows
    (follower_id, user_id)
VALUES(6, 4);


-- Insert Tweets







CREATE OR REPLACE FUNCTION addTweets
(numtimes integer)
    RETURNS VOID
AS $$
DECLARE
    text text;
BEGIN
    FOR i IN 1..numtimes LOOP
-- INSERT INTO public.users ("name", username, email, updated_at)
-- VALUES ('user_' + i, 'user', 'firstuser@gmail.com', 'now');
END
LOOP;
END;
$$
LANGUAGE plpgsql;

