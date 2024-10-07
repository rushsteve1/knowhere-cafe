CREATE OR REPLACE VIEW presentations AS
SELECT
	id,
	created_at,
	updated_at,
	deleted_at,
	users.email AS title,
	users.bio AS content,
	users.avatar AS img,
	'/users/' || users.id AS path
FROM
	users
WHERE
	deleted_at IS NULL
UNION
SELECT
	id,
	created_at,
	updated_at,
	deleted_at,
	posts.subject AS title,
	posts.body AS content,
	NULL AS img,
	'/posts/' || posts.id AS path
FROM
	posts
WHERE
	deleted_at IS NULL;
