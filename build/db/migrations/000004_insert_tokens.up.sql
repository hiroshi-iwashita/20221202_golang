INSERT INTO `tokens`
    (
        user_id,
        token,
        created_at,
        updated_at,
        expire_at
    )
VALUES
    (
        'user_id_001',
        'token_001',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        '2022-12-31 04:01:24.000'
    )
;