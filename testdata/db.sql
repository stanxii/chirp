DROP TABLE IF EXISTS users;

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
    ('Dua Lipa', 'duasings007', 'dua@lipa.com');

INSERT INTO public.users
    ("name", username, email)
VALUES
    ('Bob Dylan', 'bobbyd', 'bob@dylan.com');

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

