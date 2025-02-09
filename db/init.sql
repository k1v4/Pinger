CREATE TABLE IF NOT EXISTS containers (
                                     ip TEXT PRIMARY KEY,
                                     ping_time INTEGER NOT NULL,
                                     last_successful TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_ip ON containers (ip);