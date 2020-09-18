CREATE TABLE users (
	id          integer             PRIMARY KEY,
	balance     double precision    NOT NULL
);

CREATE TABLE transactions (
    id          integer             PRIMARY KEY,
    type        text                NOT NULL,
    from        integer,
    to          integer,
    amount      double precision    NOT NULL,
    date        integer             NOT NULL
);