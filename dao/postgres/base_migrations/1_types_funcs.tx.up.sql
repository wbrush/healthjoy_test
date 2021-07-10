CREATE
    OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TYPE template_status AS ENUM ('new', 'active', 'removed');