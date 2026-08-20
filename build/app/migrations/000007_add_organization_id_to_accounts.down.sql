ALTER TABLE auth_tokens
    DROP COLUMN organization_id;

ALTER TABLE account
    DROP CONSTRAINT account_id_organization_id_key,
    DROP COLUMN organization_id;