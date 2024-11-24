CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link VARCHAR(255),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE INDEX idx_groups_name ON groups (name);
CREATE INDEX idx_songs_group_id ON songs (group_id);
CREATE INDEX idx_songs_title ON songs (song);