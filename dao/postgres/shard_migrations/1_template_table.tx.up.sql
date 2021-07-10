CREATE TABLE ?shard.templates
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR         NOT NULL,
    template_self VARCHAR         NOT NULL,
    status        template_status NOT NULL DEFAULT 'new',
    description   TEXT            NULL     DEFAULT NULL,
    created_at    TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP       NULL,
    archived_at   TIMESTAMP       NULL
);

CREATE UNIQUE INDEX templates_name_uindex
    on ?shard.templates (name);

CREATE UNIQUE INDEX templates_self_uindex
    on ?shard.templates (template_self);

CREATE TRIGGER templates_set_updated_at BEFORE UPDATE ON ?shard.templates
    FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
