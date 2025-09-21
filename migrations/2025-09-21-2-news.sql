-- Enable trigram extension for GIN trigram indexes
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE,
    alias VARCHAR(255) NOT NULL UNIQUE,
    createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- News table
CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    alias VARCHAR(255) NOT NULL UNIQUE,
    preview TEXT,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    views INTEGER NOT NULL DEFAULT 0
);

-- News to categories many-to-many
CREATE TABLE IF NOT EXISTS news_categories (
    news_id INTEGER NOT NULL REFERENCES news(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (news_id, category_id)
);

-- Indexes to speed up search
CREATE INDEX IF NOT EXISTS idx_news_createdat ON news(createdat DESC);
CREATE INDEX IF NOT EXISTS idx_news_title_trgm ON news USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_news_content_trgm ON news USING GIN (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_categories_alias ON categories(alias);
CREATE INDEX IF NOT EXISTS idx_news_alias ON news(alias);


