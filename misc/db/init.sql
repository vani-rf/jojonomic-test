CREATE TABLE IF NOT EXISTS tbl_harga (
    reff_id VARCHAR (15) PRIMARY KEY NOT NULL,
    admin_id VARCHAR (15) NOT NULL,
    harga_topup DECIMAL(12, 2) NOT NULL,
    harga_buyback DECIMAL(12, 2) NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tbl_rekening (
    id VARCHAR (15)  PRIMARY KEY NOT NULL,
    norek VARCHAR (15) UNIQUE NOT NULL,
    customer_name VARCHAR (20) NOT NULL,
    amount DECIMAL(12, 2) DEFAULT 0,
    created_at TIMESTAMP
);

CREATE TABLE if NOT EXISTS tbl_transaksi (
    id VARCHAR (15)  PRIMARY KEY NOT NULL,
    reff_id VARCHAR (15),
    norek VARCHAR (15),
    type VARCHAR(15),
    gold_weight FLOAT
);
