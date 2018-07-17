DROP TABLE IF EXISTS users;

CREATE TABLE public.users (
    id serial NOT NULL,
    "name" text NULL,
    username text NOT NULL,
    email text NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    follower_count int4 NULL,
    following_count int4 NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id))
WITH (OIDS = FALSE);

CREATE UNIQUE INDEX uix_users_email ON public.users
USING btree (email);

CREATE UNIQUE INDEX uix_users_username ON public.users
USING btree (username);

-- Insert Users
INSERT INTO public.users ("name", username, email, created_at)
    VALUES ('Sam Smith', 'samsmith', 'sam2018@gmail.com', 'now');

INSERT INTO public.users ("name", username, email, created_at)
    VALUES ('Kanye West', 'kanye_west', 'kanye@kanye.com', 'now');

INSERT INTO public.users ("name", username, email, created_at)
    VALUES ('Dua Lipa', 'duasings007', 'dua@lipa.com', 'now');

INSERT INTO public.users ("name", username, email, created_at)
    VALUES ('Bob Dylan', 'bobbyd', 'bob@dylan.com', 'now');

CREATE OR REPLACE FUNCTION addTweets (numtimes integer)
    RETURNS VOID
AS $$
DECLARE
    text text;
BEGIN
    FOR i IN 1..numtimes LOOP
        -- INSERT INTO public.users ("name", username, email, created_at, updated_at)
        -- VALUES ('user_' + i, 'user', 'firstuser@gmail.com', 'now', 'now');
    END LOOP;
END;
$$
LANGUAGE plpgsql;

