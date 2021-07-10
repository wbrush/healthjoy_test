CREATE TABLE shard_info (
  shard_id INTEGER PRIMARY KEY,
  property_name VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NULL,
  archived_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX shard_info_uindex
    on shard_info(shard_id);

CREATE TRIGGER shard_info_set_updated_at BEFORE
    UPDATE ON shard_info FOR EACH ROW EXECUTE PROCEDURE set_updated_at();