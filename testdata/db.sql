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
    followers_count int4 NULL,
    follower_count int4 NULL,
    following_count int4 NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)
WITH (
	OIDS=FALSE
) ;
CREATE UNIQUE INDEX uix_users_email ON public.users USING btree
(email) ;
CREATE UNIQUE INDEX uix_users_remember_hash ON public.users USING btree
(remember_hash) ;
CREATE UNIQUE INDEX uix_users_username ON public.users USING btree
(username) ;

-- Insert Users
INSERT INTO public.users
    ("name", username, email)
VALUES
    ('Sam Smith', 'samsmith', 'sam2018@gmail.com');

INSERT INTO public.users
    ("name", username, email)
VALUES
    ('Kanye West', 'kanye_west', 'kanye@kanye.com');

INSERT INTO public.users
    ("name", username, email)
VALUES
    ('Dua Lipa', 'duasings', 'dua@lipa.com');

INSERT INTO public.users
    ("name", username, email)
VALUES
    ('Bob Dylan', 'bobbyd', 'bob@dylan.com');

-- this user's remember_hash based on
-- HMACKey: "secret-hmac-key" 
-- remember token: "ke3kO2KwD4HjC2lqhYWDD17T3aKXanDN1qiMLQLq1LI="
INSERT INTO public.users
    ("name", username, email, remember_hash)
VALUES
    ('Tom Tester', 'tommyTesterton', 'tommy@gmail.com', 'fe62al6IpvpqvVke8hQTgQqWyi59d1Z7ZGXSC00kTf8=');

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

