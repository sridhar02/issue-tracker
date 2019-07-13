UPDATE issues SET pinned='Unpinned' WHERE pinned IS NULL;

UPDATE issues SET lock = 'Unlocked' WHERE lock IS NULL;