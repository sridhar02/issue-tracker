CREATE TYPE status AS ENUM ('Open','Closed');

CREATE TYPE pin AS ENUM ('Pinned','Unpinned');


CREATE TYPE lock AS ENUM ('Locked','Unlocked');

ALTER TABLE issues DROP COLUMN status;

ALTER TABLE issues DROP COLUMN pinned;

ALTER TABLE issues DROP COLUMN lock;


ALTER TABLE issues ADD COLUMN status status;

ALTER TABLE issues ADD COLUMN pinned pin;

ALTER TABLE issues ADD COLUMN lock lock;