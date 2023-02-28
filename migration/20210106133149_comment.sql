-- +goose Up
CREATE TABLE comment (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  description character varying NOT NULL,
  type character varying NOT NULL,
  updated boolean NOT NULL,
  account_id character varying NOT NULL,
  advertiser_id character varying NOT NULL,
  listing_id character varying NOT NULL,
  owner JSONB NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL
);

CREATE INDEX comment_advertiser_id_account_id ON comment USING btree (advertiser_id, account_id);

-- +goose Down
DROP TABLE comment;
