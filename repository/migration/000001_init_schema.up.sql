CREATE TABLE country (
    ISO2 CHAR(2) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE bank (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL,
    address VARCHAR(200) NOT NULL,
    swiftCode CHAR(11) NOT NULL UNIQUE,
    countryISO2 CHAR(2) NOT NULL,
    isHeadquarter BOOLEAN NOT NULL,
    headquarterId INT,
    FOREIGN KEY (countryISO2)
        REFERENCES country (ISO2)
        ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (headquarterId)
        REFERENCES bank (id)
        ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE UNIQUE INDEX idxSwidftCode
ON bank (swiftCode);