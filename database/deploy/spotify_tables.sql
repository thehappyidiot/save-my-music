-- Deploy save-my-music:spotify_tables to pg
-- requires: spotify_schema
-- requires: users

BEGIN;

CREATE TABLE spotify.users
(
    spotify_id character varying(32) NOT NULL,
    display_name character varying(32),
    internal_id bigint NOT NULL,
    last_scan timestamp without time zone,
    created timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (spotify_id),
    UNIQUE (internal_id),
    FOREIGN KEY (internal_id)
        REFERENCES app.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE spotify.playlists
(
    owner character varying(32) NOT NULL,
    playlist_id character(22),
    name character varying(100) NOT NULL,
    description character varying(300),
    created timestamp with time zone NOT NULL DEFAULT NOW(),
    is_deleted boolean DEFAULT FALSE,
    PRIMARY KEY (playlist_id),
    FOREIGN KEY (owner)
        REFERENCES spotify.users (spotify_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TYPE spotify.artist AS
(
	name text,
	artist_id character(22)
);

CREATE TABLE spotify.tracks
(
    track_id character(22),
    album_id character(22) NOT NULL,
    album_name text NOT NULL,
    artists spotify.artist[] NOT NULL,
    available_markets character(2)[] NOT NULL DEFAULT '{}',
    is_explicit boolean,
    isrc character(12),
    ian numeric(13, 0),
    upc numeric(12, 0),
    restriction_reason text,
    name text NOT NULL,
    is_local boolean,
    PRIMARY KEY (track_id)
);
COMMENT ON COLUMN spotify.tracks.isrc
    IS 'https://en.wikipedia.org/wiki/International_Standard_Recording_Code';
COMMENT ON COLUMN spotify.tracks.ian
    IS 'https://en.wikipedia.org/wiki/International_Article_Number';
COMMENT ON COLUMN spotify.tracks.upc
    IS 'https://en.wikipedia.org/wiki/Universal_Product_Code';

CREATE TYPE spotify.show AS
(
	show_id character(22),
	available_markets character(2)[],
	description text,
	is_explicit boolean,
	languages character(2)[],
	name text,
	publisher text
);

CREATE TABLE spotify.episodes
(
    episode_id character(22) COLLATE pg_catalog."default" NOT NULL,
    is_explicit boolean NOT NULL,
    description text COLLATE pg_catalog."default",
    duration_ms integer NOT NULL,
    languages character(2)[] COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    release_date date NOT NULL,
    restriction_reason text COLLATE pg_catalog."default",
    show spotify.show NOT NULL,
    CONSTRAINT episodes_pkey PRIMARY KEY (episode_id)
);

CREATE TYPE spotify.item_type AS ENUM
('track', 'episode');

CREATE TABLE spotify.playlist_items
(
    playlist_id character(22) COLLATE pg_catalog."default",
    track_id character(22) COLLATE pg_catalog."default",
    item_type spotify.item_type NOT NULL,
    added_at timestamp with time zone NOT NULL,
    added_by character(32) COLLATE pg_catalog."default" NOT NULL,
    is_local boolean NOT NULL,
    is_unavailable boolean NOT NULL DEFAULT false,
    is_deleted boolean NOT NULL DEFAULT false,
    episode_id character(22) COLLATE pg_catalog."default",
    CONSTRAINT playlist_items_playlist_id_track_id_episode_id_key UNIQUE NULLS NOT DISTINCT (playlist_id, track_id, episode_id),
    CONSTRAINT playlist_items_episode_id_fkey FOREIGN KEY (episode_id)
        REFERENCES spotify.episodes (episode_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT playlist_items_playlist_id_fkey FOREIGN KEY (playlist_id)
        REFERENCES spotify.playlists (playlist_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT playlist_items_track_id_fkey FOREIGN KEY (track_id)
        REFERENCES spotify.tracks (track_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

COMMIT;
