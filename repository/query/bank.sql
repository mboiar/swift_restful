-- name: CreateBank :execresult
INSERT INTO bank(
    `address`,
    `name`,
    `countryISO2`,
    `isHeadquarter`,
    `swiftCode`
) VALUES (
    sqlc.arg('address'), sqlc.arg('bankName'), sqlc.arg('countryISO2'), sqlc.arg('isHeadquarter'), sqlc.arg('swiftCode')
);

-- name: SetBranchHeadquarter :execresult
UPDATE bank AS branch
INNER JOIN bank AS headquarter
ON LEFT(bank.swiftCode, 8) = LEFT(headquarter.swiftCode, 8)
SET
branch.headquarterId = headquarter.id
WHERE headquarter.isHeadquarter AND branch.id = ?;

-- name: UpdateBranchesHeadquarter :exec
UPDATE bank
SET
`headquarterId` = sqlc.arg('headquarterId')
WHERE LEFT(`SwiftCode`, 8) = LEFT(sqlc.arg('swiftCode'), 8) AND NOT `isHeadquarter`;

-- name: DeleteBank :exec
DELETE FROM bank
WHERE SwiftCode = sqlc.arg('swiftCode');

-- name: GetBankBySwiftCode :one
SELECT bank.*, country.name FROM bank
INNER JOIN country
ON bank.`countryISO2` = country.`ISO2`
WHERE SwiftCode = sqlc.arg('swiftCode') LIMIT 1;

-- name: GetBranchesByHeadquarterId :many
SELECT * from bank
WHERE `headquarterId` = ? LIMIT ?;

-- name: GetBranchesByCountryISO2 :many
SELECT * FROM bank
WHERE `countryISO2` = ? LIMIT ?;
