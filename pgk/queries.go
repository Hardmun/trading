package pgk

var Queries = []string{
	//*******************************************CREATE TABLES*******************************************
	`CREATE TABLE IF NOT EXISTS &table
(
	"opentime"
	INTEGER PRIMARY KEY
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
);`,
	//*******************************************FILL TRADING TABLES*******************************************
	`
	INSERT INTO &tableName (
		opentime, 
		openprice, 
		highprice, 
		lowprice, 
		closeprice, 
		volume, 
		closetime, 
		quoteassetvolume, 
		tradesnumber, 
		takerbaseasset, 
		takerquoteasset
	) 
	VALUES(
		?,?,?,?,?,?,?,?,?,?,?
	)                     
	ON CONFLICT (opentime) 
	DO NOTHING`,
	//*******************************************GET MAXIMUM PERIOD FOR TABLE**************************************
	`SELECT
	max(closetime)
	FROM &tableName`,
}
