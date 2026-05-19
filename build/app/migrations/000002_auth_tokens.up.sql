-- CREATE TYPE IF NOT EXISTS token_type AS ENUM('access', 'refresh')

-- CREATE TABLE IF NOT EXISTS auth_tokens (
--     id               BIGSERIAL PRIMARY KEY,         
--     user_id          UUID NOT NULL,                 
--     token_type       token_type,          
--     token_hash       TEXT NOT NULL,   
--     expires_at       TIMESTAMP NOT NULL,             
--     revoked          BOOLEAN DEFAULT FALSE,          
--     created_at       TIMESTAMP DEFAULT NOW(),
--     updated_at       TIMESTAMP DEFAULT NOW(),

--     ip_address       INET,                           
--     user_agent       TEXT,                           
--     last_used_at     TIMESTAMP,                      
--     parent_token_id  BIGINT,                        

--     FOREIGN KEY (user_id) REFERENCES account(id) ON DELETE CASCADE,
--     FOREIGN KEY (parent_token_id) REFERENCES auth_tokens(id) ON DELETE SET NULL
-- );