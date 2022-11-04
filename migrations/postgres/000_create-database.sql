DROP TABLE IF EXISTS card;
DROP TABLE IF EXISTS card_user_progress;
DROP TABLE IF EXISTS card_metrics;
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS collection;
DROP TABLE IF EXISTS collection_user_progress;
DROP TABLE IF EXISTS collection_user_metrics;
DROP TABLE IF EXISTS collection_metrics;
DROP TABLE IF EXISTS collection_cards;
DROP TYPE IF EXISTS card_user_progress_status_enum;

CREATE TABLE collection_cards (
    id uuid NOT NULL,
    card_id uuid NOT NULL,
    collection_id uuid NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE card (
    id uuid NOT NULL,
    word VARCHAR (150) NOT NULL,
    image_url TEXT NOT NULL,
    definition TEXT NOT NULL,
    sentence TEXT NOT NULL,
    antonyms VARCHAR (150) NULL,
    synonyms VARCHAR (150) NULL,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE TYPE card_user_progress_status_enum AS enum('mastered','reviewing','learning', 'none');

CREATE TABLE card_user_progress (
    id uuid NOT NULL,
    card_id uuid NOT NULL,
    user_id uuid NOT NULL,
    status card_user_progress_status_enum NOT NULL default 'none',
    learning_count INT NOT NULL default 0,
    PRIMARY KEY (id)
);

CREATE TABLE card_metrics (
    id uuid NOT NULL,
    card_id uuid NOT NULL,
    likes INT NOT NULL default 0,
    dislikes INT NOT NULL default 0,
    PRIMARY KEY (id)
);

CREATE TABLE "user" (
    id uuid NOT NULL,
    name VARCHAR (150) NOT NULL,
    email VARCHAR (150) NOT NULL,
    password TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE collection (
    id uuid NOT NULL,
    name VARCHAR (150) NOT NULL,
    author_id uuid NOT NULL,
    topics TEXT[] NOT NULL DEFAULT array[]::VARCHAR[],
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE collection_user_progress (
    id uuid NOT NULL,
    collection_id uuid NOT NULL,
    user_id uuid NOT NULL,
    mastered INT NOT NULL default 0,
    reviewing INT NOT NULL default 0,
    learning INT NOT NULL default 0,
    PRIMARY KEY (id)
);
CREATE TABLE collection_user_metrics (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    collection_id uuid NOT NULL,
    liked boolean NOT NULL default FALSE,
    disliked boolean NOT NULL default FALSE,
    viewed boolean NOT NULL default FALSE,
    starred boolean NOT NULL default FALSE,
    PRIMARY KEY (id)
);
CREATE TABLE collection_metrics (
    id uuid NOT NULL,
    collection_id uuid NOT NULL,
    likes INT NOT NULL default 0,
    dislikes INT NOT NULL default 0,
    views INT NOT NULL default 0,
    PRIMARY KEY (id)
);