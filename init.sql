CREATE TABLE users (
	id          integer             PRIMARY KEY,
	balance     double precision    NOT NULL
);

CREATE TABLE transactions (
    id          SERIAL              PRIMARY KEY,
    type        text                NOT NULL,
    sender      integer             NULL,
    receiver    integer             NULL,
    amount      double precision    NOT NULL,
    time        integer             NOT NULL
);