CREATE TABLE IF NOT EXISTS characters
(
    id          bigserial PRIMARY KEY,
    name        text                        NOT NULL,
    age         int                         NOT NULL,
    gender      text                        NOT NULL,
    status      text                        NOT NULL,
    nation      text                        NOT NULL,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at  timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS episodes
(
    id              bigserial PRIMARY KEY,
    title           text                        NOT NULL,
    air_date        date                        NOT NULL,
    created_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at      timestamp(0) with time zone NOT NULL DEFAULT NOW(),
);

CREATE TABLE IF NOT EXISTS characters_and_episodes
(
    id         bigserial PRIMARY KEY,
    character_id     bigserial                        NOT NULL,
    episode_id       bigserial                        NOT NULL,
    created_at       timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at       timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    FOREIGN KEY (character_id)
        REFERENCES characters(id),
    FOREIGN KEY (episode_id)
        REFERENCES episodes(id)
);
