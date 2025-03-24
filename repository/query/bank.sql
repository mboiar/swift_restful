-- name: CreateBank :execresult
INSERT INTO bank(
    `address`,
    `name`,
    `country_ISO2`,
    `is_headquarter`,
    `swift_code`
) VALUES (
    sqlc.arg('address'), sqlc.arg('bankName'), sqlc.arg('countryISO2'), sqlc.arg('isHeadquarter'), sqlc.arg('swiftCode')
);

-- name: SetBranchHeadquarter :execresult
UPDATE bank AS branch
INNER JOIN bank AS headquarter
ON LEFT(bank.swift_code, 8) = LEFT(headquarter.swift_code, 8)
SET
branch.headquarter_id = headquarter.id
WHERE headquarter.is_headquarter AND branch.id = ?;

-- name: UpdateBranchesHeadquarter :exec
UPDATE bank
SET
`headquarter_id` = sqlc.arg('headquarterId')
WHERE LEFT(`swift_code`, 8) = LEFT(sqlc.arg('swiftCode'), 8) AND NOT `is_headquarter`;

-- name: DeleteBank :exec
DELETE FROM bank
WHERE swift_code = sqlc.arg('swiftCode');

-- name: GetBankBySwiftCode :one
SELECT bank.*, country.name FROM bank
INNER JOIN country
ON bank.`country_ISO2` = country.`ISO2`
WHERE swift_code = sqlc.arg('swiftCode') LIMIT 1;

-- name: GetBranchesByHeadquarterId :many
SELECT * from bank
WHERE `headquarter_id` = ? LIMIT ?;

-- name: GetBranchesByCountryISO2 :many
SELECT * FROM bank
WHERE `country_ISO2` = ? LIMIT ?;
