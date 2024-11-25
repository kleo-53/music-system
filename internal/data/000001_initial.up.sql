create table if not exists songs(
    id serial primary key,
    song_group varchar not null,
    song varchar not null,
    song_text varchar,
    release_date varchar,
    link varchar
)