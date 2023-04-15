CREATE TABLE transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    account TEXT NOT NULL,
    booking_date TEXT NOT NULL,
    valuta_date TEXT NOT NULL,
    booking_text TEXT,
    purpose TEXT,
    creditor_id TEXT,
    mandate_ref TEXT,
    customer_ref TEXT,
    collector_ref TEXT,
    orig_amount REAL,
    chargeback_fee REAL,
    beneficiary TEXT,
    account_number TEXT,
    bic TEXT,
    amount REAL NOT NULL,
    currency TEXT NOT NULL,
    additional_details TEXT
);

