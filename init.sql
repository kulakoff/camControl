CREATE TABLE IF NOT EXISTS cameras (
    id SERIAL PRIMARY KEY,
    ip VARCHAR(15) NOT NULL,
    login VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL
);

-- test cam
INSERT INTO cameras (ip, login, password)
VALUES ('192.168.13.100', 'admin', 'password')
RETURNING id;
