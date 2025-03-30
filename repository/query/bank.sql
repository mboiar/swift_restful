-- name: CreateBank :execresult
INSERT INTO bank(
    `address`,
    `name`,
    `country_ISO2`,
    `swift_code`
) VALUES (NULLIF(?, ''), ?, ?, ?);

-- name: CreateBankBulk :copyfrom
INSERT INTO bank(
    `address`,
    `name`,
    `country_ISO2`,
    `swift_code`
) VALUES (
    ?, ?, ?, ?
);

-- name: DeleteBank :exec
DELETE FROM bank
WHERE swift_code = ?;

-- name: GetBankBySwiftCode :one
SELECT bank.*, country.name FROM bank
INNER JOIN country
ON bank.`country_ISO2` = country.`ISO2`
WHERE swift_code = ? LIMIT 1;

-- name: GetBranchesBySwiftCode :many
SELECT * from bank
WHERE LEFT(bank.swift_code, 8) = LEFT(?, 8) LIMIT ?;

-- name: GetBranchesByCountryISO2 :many
SELECT * FROM bank
WHERE `country_ISO2` = ? LIMIT ?;
