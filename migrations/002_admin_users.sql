CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    pass_hash TEXT NOT NULL,
    created_at timestamptz DEFAULT now()
);

-- initial admin example pw=123456
INSERT INTO admin_users(email, pass_hash)
VALUES (
           'admin@example.com',
           '$2a$12$tlcJH6uawSs4u.0Tic9m7uU0hiuHV0yaf3KqN0bMoA5vgxJda7Mhy'
       ) ON CONFLICT DO NOTHING;