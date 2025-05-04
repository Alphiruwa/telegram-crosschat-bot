CREATE TABLE IF NOT EXISTS links
(
    src_chat_id BIGINT,
    tgt_chat_id BIGINT,
    created_at DATE NOT NULL DEFAULT NOW(),   
    PRIMARY KEY(src_chat_id, tgt_chat_id)
);

CREATE TABLE IF NOT EXISTS requests
(
    src_chat_id BIGINT,
    tgt_chat_id BIGINT,
    created_at DATE NOT NULL DEFAULT NOW(),
    tgt_message_id BIGINT,
    from_user_id BIGINT, 
    PRIMARY KEY(src_chat_id, tgt_chat_id)
);