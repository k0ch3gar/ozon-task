WITH user_id AS (
    INSERT INTO
        users (username, email, password)
    VALUES
        ('foo', 'foo@mail.ru', '1')
    RETURNING id
)
INSERT INTO posts (author_id, title, body, allow_comments) VALUES ((SELECT id FROM user_id),'foo', 'foofoofoo', true);

WITH user_id AS (
    INSERT INTO
        users (username, email, password)
        VALUES
            ('bar', 'bar@mail.ru', '2')
        RETURNING id
)
INSERT INTO posts (author_id, title, body, allow_comments) VALUES ((SELECT id FROM user_id),'bar', 'barbarbar', false);

WITH user_id AS (
    INSERT INTO
        users (username, email, password)
        VALUES
            ('baz', 'baz@mail.ru', '3')
        RETURNING id
)
INSERT INTO posts (author_id, title, body, allow_comments) VALUES ((SELECT id FROM user_id),'baz', 'bazbazbaz', true);