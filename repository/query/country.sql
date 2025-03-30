-- name: CreateCountry :execresult
INSERT IGNORE INTO country(
    `ISO2`,
    name
) VALUES (?, ?);

-- name: CreateCountryBulk :copyfrom
INSERT IGNORE INTO country(
    `ISO2`,
    name
) VALUES (
    ?, ?
);

-- name: GetCountryByCountryISO2 :one
SELECT * FROM country
WHERE `ISO2` = ? LIMIT 1;