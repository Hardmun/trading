CREATE TABLE IF NOT EXISTS &table
(
    "opentime"
        INTEGER
        NOT NULL,
    "openprice"
        REAL
        NOT NULL,
    "highprice"
        REAL
        NOT NULL,
    "lowprice"
        REAL
        NOT NULL ,
    "closeprice"
        REAL
        NOT NULL,
    "volume"
        REAL
        NOT NULL,
    "closetime"
        INTEGER
        NOT NULL,
    "quoteassetvolume"
        REAL
        NOT NULL,
    "tradesnumber"
        INTEGER
        NOT NULL,
    "takerbaseasset"
        REAL
        NOT NULL,
    "takerquoteasset"
        REAL
        NOT NULL
);
