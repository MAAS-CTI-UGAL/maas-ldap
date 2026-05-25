CREATE TABLE IF NOT EXISTS maas_user_mappings (
    username TEXT PRIMARY KEY,
    maas_password TEXT NOT NULL
);

INSERT INTO maas_user_mappings (username, maas_password)
VALUES
    ('student1', 'parola1'),
    ('tb171', 'parola1234')
ON CONFLICT(username) DO UPDATE SET
    maas_password = excluded.maas_password;
