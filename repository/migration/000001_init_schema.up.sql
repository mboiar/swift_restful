CREATE TABLE country (
    ISO2 CHAR(2) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE bank (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL,
    address VARCHAR(200) NOT NULL,
    swift_code CHAR(11) NOT NULL UNIQUE,
    country_ISO2 CHAR(2) NOT NULL,
    is_headquarter BOOLEAN NOT NULL,
    headquarter_id INT,
    FOREIGN KEY (country_ISO2)
        REFERENCES country (ISO2)
        ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (headquarter_id)
        REFERENCES bank (id)
        ON DELETE SET NULL ON UPDATE CASCADE
);