DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS tweets;
DROP TABLE IF EXISTS likes;

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

CREATE TABLE public.tweets
(
    id serial NOT NULL,
    post text NULL,
    username text NOT NULL,
    likes_count int4 NULL,
    retweets_count int4 NULL,
    retweet_id int4 NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT tweets_pkey PRIMARY KEY (id)
)
WITH (
	OIDS=FALSE
) ;
CREATE INDEX idx_tweets_deleted_at ON public.tweets USING btree
(deleted_at) ;
CREATE INDEX idx_tweets_username ON public.tweets USING btree
(username) ;

CREATE TABLE public.likes
(
    tweet_id serial NOT NULL,
    user_id serial NOT NULL,
    CONSTRAINT likes_pkey PRIMARY KEY (tweet_id, user_id)
)
WITH (
	OIDS=FALSE
) ;



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
    ('Tom Tester', 'tommytesterton', 'tommy@gmail.com', 'fake-pw-hash',
        'M-QQp9qW_mhAZK5s1651k_ODELJZTP1EBahOOByjfJE=');
-- this user's remember_hash based on
-- HMACKey: "secret-hmac-key" 
INSERT INTO public.users
    ("name", username, email, password_hash, remember_hash)
VALUES
    ('Vince Main', 'vincetester', 'vtester@gmail.com', 'fake-pw-hash',
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


-- Insert 25 tweets
INSERT INTO public.tweets
    (id, username, post, retweet_id)
VALUES(1001, 'duasings', 'Hey, this is my first tweet!', 0);
INSERT INTO public.tweets
    (id, username, post, retweet_id)
VALUES(1002, 'duasings', 'Second tweet! Let''s go!', 0);
INSERT INTO public.tweets
    (id, username, post, retweet_id)
VALUES(1003, 'bobbyd', 'I love playing the guitar.', 0);
INSERT INTO public.tweets
    (id, username, post, retweet_id)
VALUES(1004, 'vincetester', 'this tweet will be deleted...', 0);
INSERT INTO public.tweets
    (id, username, post, retweet_id)
VALUES(1005, 'vincetester', 'this tweet will be updated...', 0);
INSERT INTO public.tweets
    (id, username, post, retweet_id)
VALUES(1006, 'kanye_west', 'amazing tweet by kanye', 0);
-- INSERT INTO public.tweets
--     (id, username, post, retweet_id)
-- VALUES(1007, 'bobbyd', 'ma', 0);
INSERT INTO public.likes
    (tweet_id, user_id)
VALUES(1003, 1);
INSERT INTO public.likes
    (tweet_id, user_id)
VALUES(1003, 2);
INSERT INTO public.likes
    (tweet_id, user_id)
VALUES(1003, 3);
INSERT INTO public.likes
    (tweet_id, user_id)
VALUES(1006, 6);




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

