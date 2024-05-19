-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO "snippets" ("id", "title", "content", "create_time", "expire_time")
VALUES ('334d7468-f258-4f69-b5e0-f3ff6f265c75',
        'An old silent pond',
        'An old silent pond...' || char(10) || 'A frog jumps into the pond,' || char(10) || 'splash! Silence again.' ||
        char(10) || char(10) || '– Matsuo Bashō',
        current_timestamp,
        datetime(current_timestamp, '365 days'));

INSERT INTO "snippets" ("id", "title", "content", "create_time", "expire_time")
VALUES ('a1483631-54f0-4401-82b7-4b86406570b5',
        'Over the wintry forest',
        'Over the wintry' || char(10) || 'forest, winds howl in rage' || char(10) || 'with no leaves to blow.' ||
        char(10) || char(10) || '– Natsume Soseki',
        current_timestamp,
        datetime(current_timestamp, '365 days'));

INSERT INTO "snippets" ("id", "title", "content", "create_time", "expire_time")
VALUES ('092165b6-4165-4676-ab2d-354dc5bb712b',
        'First autumn morning',
        'First autumn morning' || char(10) || 'the mirror I stare into' || char(10) || 'shows my father''s face.' ||
        char(10) || char(10) || '– Murakami Kijo',
        current_timestamp,
        datetime(current_timestamp, '7 days'));
