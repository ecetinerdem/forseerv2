CREATE TABLE IF NOT EXISTS portfolios(
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    name varchar(50) UNIQUE NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);