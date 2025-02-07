CREATE TABLE IF NOT EXISTS containers (
                                     ip TEXT PRIMARY KEY,
                                     ping_time DATE NOT NULL,
                                     last_successful DATE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_ip ON containers (ip);