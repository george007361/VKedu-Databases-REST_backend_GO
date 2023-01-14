-- CITEXT - тот же TEXT только без учета регистра, авто lower при сравнении
-- UNLOGGED - нежурналируемые. ВЫигрываем по времени, теряем по аварийно-безопасности
-- COLLATE "C" - сравнение по порядку байт

-- Подключаем citext
CREATE EXTENSION IF NOT EXISTS CITEXT;

-- Чистим бд
DROP TABLE IF EXISTS threads,
forums,
users,
votes,
posts,
nickname_forum CASCADE;


-- Таблица пользователей
CREATE UNLOGGED TABLE users (
    nickname CITEXT COLLATE "C" NOT NULL PRIMARY KEY,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT UNIQUE NOT NULL
);


-- Таблица форумов
CREATE UNLOGGED TABLE forums (
    slug CITEXT NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    author_nickname CITEXT REFERENCES users(nickname) NOT NULL,
    posts BIGINT DEFAULT 0,
    threads BIGINT DEFAULT 0
);

-- Таблица постов
CREATE UNLOGGED TABLE posts (
    id BIGSERIAL  PRIMARY KEY,
    parent_id BIGINT DEFAULT 0,
    author_nickname CITEXT  REFERENCES users(nickname),
    message TEXT,
    isedited BOOLEAN DEFAULT false,
    forum_slug CITEXT  REFERENCES forums(slug),
    thread_id BIGINT,
    created timestamp with time zone ,
    path_tree BIGINT[]
);

-- Таблица веток
CREATE UNLOGGED TABLE threads (
    id BIGSERIAL  PRIMARY KEY,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    slug CITEXT UNIQUE,
    forum_slug CITEXT  REFERENCES forums(slug) NOT NULL,
    author_nickname CITEXT  REFERENCES users(nickname) NOT NULL,
    created timestamp with time zone,
    votes BIGINT DEFAULT 0
);

-- Таблица голосов
CREATE UNLOGGED TABLE votes (
    id BIGSERIAL  PRIMARY KEY,
    nickname CITEXT  REFERENCES users(nickname),
    voice BIGINT,
    thread_id BIGINT NOT NULL  REFERENCES threads(id),
    UNIQUE (thread_id, nickname)
);

-- Таблица пользователей формума
CREATE UNLOGGED TABLE forum_users (
    nickname CITEXT COLLATE "C" REFERENCES users(nickname),
    fullname TEXT,
    about TEXT,
    email CITEXT,
    forum_slug CITEXT REFERENCES forums(slug),
    UNIQUE (forum_slug, nickname)
);


--  Обновление количества веток на форуме при создании новой
CREATE OR REPLACE FUNCTION func_after_thread_CREATEd_update_threads() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums SET threads = threads + 1
    WHERE slug = NEW.forum_slug;
RETURN NULL;
END
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER trig_thread_added_count_threads
AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE func_after_thread_created_update_threads();

--  Обновление количества постов на форуме при создании нового
CREATE OR REPLACE FUNCTION func_after_post_CREATEd_update_forum_posts() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums SET posts = posts + 1
    WHERE slug = NEW.forum_slug;
RETURN NULL;
END
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER trig_post_added_count_posts
AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE func_after_post_created_update_forum_posts();

-- Добавление пользователя к форуму при создании ветки
CREATE OR REPLACE FUNCTION func_after_thread_CREATEd_add_user_to_forum() RETURNS TRIGGER AS $$
DECLARE
    author_fullname TEXT;
    author_about TEXT;
    author_email CITEXT;
BEGIN
    SELECT fullname, about, email FROM users
    WHERE nickname = NEW.author_nickname 
    INTO author_fullname, author_about, author_email;

    INSERT INTO forum_users (
        nickname,
        fullname,
        about,
        email,
        forum_slug 
    ) values (
        NEW.author_nickname,
        author_fullname,
        author_about,
        author_email,
        NEW.forum_slug
    )
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER trig_thread_added_add_user_to_forum
AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE func_after_thread_created_add_user_to_forum();

-- Добавление пользователя к форуму при создании поста
CREATE TRIGGER trig_post_added_add_user_to_forum
AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE func_after_thread_created_add_user_to_forum();

-- Обновление пути к посту
CREATE OR REPLACE FUNCTION func_update_post_path()
RETURNS TRIGGER AS $$
BEGIN
    NEW.path_tree = (SELECT path_tree FROM posts WHERE id=NEW.parent_id) || NEW.id;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trig_insert_post
before INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE func_update_post_path();

-- Обновление рейтинга thread после создания vote
CREATE OR REPLACE FUNCTION func_update_thread_votes_after_insert()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE threads SET votes = votes + NEW.voice WHERE id = NEW.thread_id;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trig_after_insent_vote
AFTER INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE func_update_thread_votes_after_insert();

-- Обновление рейтинга thread после изменения vote
CREATE OR REPLACE FUNCTION func_update_thread_votes_after_update()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE threads SET votes = votes + NEW.voice - old.voice WHERE id = NEW.thread_id;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER trig_after_update_vote
AFTER UPDATE ON votes FOR EACH ROW EXECUTE PROCEDURE func_update_thread_votes_after_update();


-- Индексы
CREATE INDEX IF NOT EXISTS index_forums_user_nickname ON forums (author_nickname);

CREATE INDEX IF NOT EXISTS index_threads_author ON threads (author_nickname);
CREATE INDEX IF NOT EXISTS index_threads_forum ON threads (forum_slug);
CREATE INDEX IF NOT EXISTS index_threads_slug ON threads (slug);
CREATE INDEX IF NOT EXISTS index_threads_forum_CREATEd ON threads (forum_slug, created);

CREATE INDEX IF NOT EXISTS index_forum_user_forum_user_nickname ON forum_users (forum_slug, nickname);

CREATE INDEX IF NOT EXISTS index_posts_thread_id ON posts (thread_id, id);
CREATE INDEX IF NOT EXISTS index_posts_thread_post_tree ON posts (thread_id, path_tree);
CREATE INDEX IF NOT EXISTS index_posts_parent_thread_id ON posts (parent_id, thread_id, id);
CREATE INDEX IF NOT EXISTS index_posts_post_tree_one_post_tree ON posts ((path_tree[1]), path_tree);

CREATE INDEX IF NOT EXISTS index_users_email ON users (email);
CREATE INDEX IF NOT EXISTS index_users_email_nickname ON users (email, nickname);

VACUUM ANALYZE;