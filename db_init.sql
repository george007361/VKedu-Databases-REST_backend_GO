-- citext - тот же TEXT только без учета регистра, авто lower при сравнении
create extension if not exists citext;

-- Чистим бд
drop table if exists threads,
forums,
users,
votes,
posts,
nickname_forum cascade;

-- unlogged - нежурналируемые. ВЫигрываем по времени, теряем по аварийно-безопасности
-- collate "C" - сравнение по порядку байт
create unlogged table users (
    nickname citext collate "C" constraint users_pk not null primary key,
    fullname text not null,
    about text,
    email citext unique not null
);

create unlogged table forums (
    slug citext not null constraint forums_pk primary key,
    title text not null,
    "user" citext constraint forums_user_fk references users(nickname) not null,
    posts bigint default 0,
    threads bigint default 0
);

--     path bigint [] default array [] :: bigint []
create unlogged table posts (
    id bigserial constraint posts_pk primary key,
    parent bigint default 0,
    author citext constraint posts_author_fk references users(nickname),
    message text,
    isedited boolean default false,
    forum citext constraint posts_forum_fk references forums(slug),
    thread integer,
    created timestamp with time zone default CURRENT_TIMESTAMP,
    path integer[]
);

create unlogged table threads (
    id bigserial constraint threads_pk primary key,
    created timestamp with time zone default CURRENT_TIMESTAMP,
    votes integer default 0,
    
    title text not null,
    author citext constraint threads_author_fk references users(nickname) not null,
    forum citext constraint threads_forum_fk references forums(slug) not null,
    message text not null,
    slug citext unique
);

create unlogged table votes (
    id bigserial constraint votes_pk primary key,
    nickname citext constraint votes_nickname_fk references users(nickname),
    voice integer,
    thread integer not null constraint votes_thread_fk references threads(id),
    unique (thread, nickname)
);

create unlogged table nickname_forum (
    nickname citext collate "C" references users(nickname),
    fullname text,
    about text,
    email citext,
    forum citext references forums,
    unique (forum, nickname)
);


--  Обновление количества веток на форуме при создании новой
create or replace function func_after_thread_created_update_threads() returns trigger as $$
begin
    update forums set threads = threads + 1
    where new.forum = slug;
return null;
end
$$ language 'plpgsql';

create trigger trig_thread_added_count_threads
after insert on threads for each row execute procedure func_after_thread_created_update_threads();

--  Обновление количества постов на форуме при создании нового
create or replace function func_after_post_created_update_forum_posts() returns trigger as $$
begin
    update forums set posts = posts + 1
    where new.forum = slug;
return null;
end
$$ language 'plpgsql';

create trigger trig_post_added_count_posts
after insert on posts for each row execute procedure func_after_post_created_update_forum_posts();



-- Добавление пользователя к форуму при создании ветки
create or replace function func_after_thread_created_add_user_to_forum() returns trigger as $$
declare
    author_fullname text;
    author_about text;
    author_email citext;
begin
    select fullname, about, email from users
    where nickname = new.author 
    into author_fullname, author_about, author_email;

    insert into nickname_forum (
        nickname,
        fullname,
        about,
        email,
        forum 
    ) values (
        new.author,
        author_fullname,
        author_about,
        author_email,
        new.forum
    )
    on conflict do nothing;
    return new;
end
$$ language 'plpgsql';

create trigger trig_thread_added_add_user_to_forum
after insert on threads for each row execute procedure func_after_thread_created_add_user_to_forum();


-- Обновление пути к посту
create or replace function func_update_post_path()
returns trigger as $$
begin
    new.path = (select path from posts where id=new.parent) || new.id;
    return new;
end
$$ LANGUAGE plpgsql;

create trigger trig_insert_post
before insert on posts for each row execute procedure func_update_post_path();


-- Обновление рейтинга thread после создания vote
create or replace function func_update_thread_votes_after_insert()
returns trigger as $$
begin
    update threads set votes = votes + new.voice where id = new.thread;
    return new;
end
$$ language plpgsql;

create trigger trig_after_insent_vote
after insert on votes for each row execute procedure func_update_thread_votes_after_insert();

-- Обновление рейтинга thread после изменения vote
create or replace function func_update_thread_votes_after_update()
returns trigger as $$
begin
    update threads set votes = votes + new.voice - old.voice where id = new.thread;
    return new;
end
$$ language plpgsql;

create trigger trig_after_update_vote
after update on votes for each row execute procedure func_update_thread_votes_after_update();





    -- create
    -- or replace function add_post() returns trigger as $ $ begin
    -- update
    --     forum
    -- set
    --     posts = posts + 1
    -- where
    --     slug = new.forum;
    -- new.path = (
    --     select
    --         path
    --     from
    --         post
    --     where
    --         id = new.parent
    --     limit
    --         1
    -- ) || new.id;
    -- return new;
    -- end $ $ language 'plpgsql';


    -- create trigger add_post before
    -- insert
    --     on post for each row execute procedure add_post();


    -- create
    -- or replace function add_vote() returns trigger as $ $ begin
    -- update
    --     thread
    -- set
    --     votes =(votes + new.voice)
    -- where
    --     id = new.thread;
    -- return new;
    -- end $ $ language 'plpgsql';
    -- create trigger add_vote
    -- after
    -- insert
    --     on votes for each row execute procedure add_vote();
    -- create
    -- or replace function update_vote() returns trigger as $ $ begin if old.voice <> new.voice then
    -- update
    --     thread
    -- set
    --     votes = votes - old.voice + new.voice
    -- where
    --     id = new.thread;
    -- end if;
    -- return new;
    -- end $ $ language 'plpgsql';
    -- create trigger update_vote
    -- after
    -- update
    --     on votes for each row execute procedure update_vote();
    -- create
    -- or replace function add_post_user() returns trigger as $ $ declare author_nickname citext;
    -- author_fullname text;
    -- author_about text;
    -- author_email citext;
    -- begin
    -- select
    --     nickname,
    --     fullname,
    --     about,
    --     email
    -- from
    --     users
    -- where
    --     nickname = new.author into author_nickname,
    --     author_fullname,
    --     author_about,
    --     author_email;
    -- insert into
    --     nickname_forum (nickname, fullname, about, email, forum)
    -- values
    --     (
    --         author_nickname,
    --         author_fullname,
    --         author_about,
    --         author_email,
    --         new.forum
    --     ) on conflict do nothing;
    -- return new;
    -- end $ $ language 'plpgsql';
    -- create trigger add_post_user
    -- after
    -- insert
    --     on post for each row execute procedure add_post_user();
    

    
    -- create index if not exists for_user_nickname on users using hash (nickname);
    -- create index if not exists for_user_email on users using hash (email);
    -- create index if not exists for_forum_slug on forum using hash (slug);
    -- create index if not exists for_thread_slug on thread using hash (slug);
    -- create index if not exists for_thread_forum on thread using hash (forum);
    -- create index if not exists for_post_thread on post using hash (thread);
    -- create index if not exists for_thread_created on thread (created);
    -- create index if not exists for_thread_created_forum on thread (forum, created);
    -- create index if not exists for_post_path_single on post ((path [1]));
    -- create index if not exists for_post_id_path_single on post (id, (path [1]));
    -- create index if not exists for_post_path on post (path);
    -- create unique index if not exists for_votes_nickname_thread_nickname on votes (thread, nickname);
    -- create index if not exists for_nickname_forum on nickname_forum using hash (nickname);
    -- create index if not exists for_nickname_forum_nickname on nickname_forum (forum, nickname);
    -- vacuum;
    -- vacuum analyze;
    -- -- -- ОЧИСТКА
    -- -- DROP TABLE IF EXISTS votes CASCADE;
    -- -- DROP TABLE IF EXISTS users CASCADE;
    -- -- DROP TABLE IF EXISTS posts;
    -- -- DROP TABLE IF EXISTS forums CASCADE;
    -- -- DROP TABLE IF EXISTS threads CASCADE;
    -- -- DROP TABLE IF EXISTS forum_users CASCADE;
    -- -- -- сОЗДАНИЕ ТАБЛИЦ
    -- -- CREATE TABLE users (
    -- --     nickname TEXT PRIMARY KEY NOT NULL UNIQUE,
    -- --     fullname TEXT,
    -- --     about TEXT,
    -- --     email TEXT UNIQUE
    -- -- );




-- \echo "  --Done" 
-- \echo "  -- -- Done" 
--  \c postgres 
--  \dt

-- insert into users(nickname, fullname, about, email)
-- values('George1', 'George Illar', 'Hello!!!', 'email1@mail.ru'),
-- ('George2', 'George Illar', 'Hello!!!', 'email2@mail.ru'),
-- ('George3', 'George Illar', 'Hello!!!', 'email3@mail.ru');
    

-- select * from users;
    
-- insert into forums (
--     slug,
--     title,
--     "user",
--     posts,
--     threads
-- )
-- values (
--     'Forum1', 
--     'Title1', 
--     (SELECT nickname
--          FROM users
--         WHERE nickname = 'George1'),
--         123,
--         321
-- );

-- select * from forums;

-- insert into threads (
--     title,
--     author,
--     forum,
--     message,
--     slug
-- )
-- values (
--     'Thread1',
--     (select nickname from users where nickname = 'george2'),
--     (select slug from forums where slug = 'Forum1'),
--     'Msg1',
--     'Slug11'
-- );


-- select * from threads;

-- insert into posts(
--     author,
--     message,
--     forum,
--     thread
-- )
-- values(
--     (select nickname from users where nickname = 'George1'),
--     'FIRST POST!!!',
--     (select slug from forums where slug = 'forum1'),
--     (select id from threads where title = 'Thread1')
-- );

-- select * from posts;

-- \q