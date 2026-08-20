ALTER TABLE account
    ADD COLUMN organization_id BIGINT NOT NULL REFERENCES organizations(id),
    ADD CONSTRAINT account_id_organization_id_key UNIQUE (id, organization_id);

ALTER TABLE auth_tokens
    ADD COLUMN organization_id BIGINT NOT NULL REFERENCES organizations(id);