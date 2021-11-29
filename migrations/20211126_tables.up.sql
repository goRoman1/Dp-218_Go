CREATE TABLE IF NOT EXISTS roles
(
    id         smallint PRIMARY KEY,
    name       VARCHAR(50) UNIQUE,
    is_admin    boolean,
    is_user     boolean,
    is_supplier boolean
);

CREATE TABLE IF NOT EXISTS users
(
    id          serial PRIMARY KEY,
    login_email  VARCHAR(100) UNIQUE NOT NULL,
    is_blocked   boolean,
    user_name    VARCHAR(100),
    user_surname VARCHAR(100),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role_id      int                 NOT NULL,

    FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE TABLE IF NOT EXISTS login_info
(
    user_id       int PRIMARY KEY,
    password_hash VARCHAR(512),

    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS login_status
(
    user_id    int PRIMARY KEY,
    logged_in  boolean,
    date_time  TIMESTAMP NOT NULL,
    ip_address VARCHAR(40),

    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS contact_types
(
    id   smallserial PRIMARY KEY,
    name VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS contacts
(
    id          serial PRIMARY KEY,
    type_id      int NOT NULL,
    user_id      int NOT NULL,
    contact_info VARCHAR(200),

    FOREIGN KEY (type_id) REFERENCES contact_types (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS accounts
(
    id       serial PRIMARY KEY,
    name     VARCHAR(100),
    number   VARCHAR(100) UNIQUE NOT NULL,
    owner_id int                 NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS payment_types
(
    ID   smallserial PRIMARY KEY,
    Name VARCHAR(100) UNIQUE
);

CREATE TABLE IF NOT EXISTS supplier_commissions
(
    id                  serial PRIMARY KEY,
    commission_percent  NUMERIC(4, 2),
    user_id             int NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS supplier_prices
(
    id               serial PRIMARY KEY,
    price            NUMERIC(15, 2),
    payment_type_id  smallint NOT NULL,
    user_id          int      NOT NULL,

    FOREIGN KEY (payment_type_id) REFERENCES payment_types (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS scooter_models
(
    id               smallserial PRIMARY KEY,
    payment_type_id  smallint     NOT NULL,
    model_name       VARCHAR(100) NOT NULL,
    max_weight       NUMERIC(5, 2),
    speed            smallint     NOT NULL,

    FOREIGN KEY (payment_type_id) REFERENCES payment_types (id)
);

CREATE TABLE IF NOT EXISTS scooters
(
    id            serial PRIMARY KEY,
    model_id      smallint            NOT NULL,
    owner_id      int                 NOT NULL,
    serial_number VARCHAR(100) UNIQUE NOT NULL,

    FOREIGN KEY (model_id) REFERENCES scooter_models (id),
    FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS locations
(
    id        serial PRIMARY KEY,
    latitude  NUMERIC(10, 0) NOT NULL,
    longitude NUMERIC(10, 0) NOT NULL,
    label     VARCHAR(200)
);

CREATE TABLE IF NOT EXISTS scooter_stations
(
    id          serial PRIMARY KEY,
    location_id int NOT NULL,
    name        VARCHAR(100),
    is_active   boolean,

    FOREIGN KEY (location_id) REFERENCES locations (id)
);

CREATE TABLE IF NOT EXISTS scooter_statuses
(
    scooter_id     int PRIMARY KEY,
    location_id    int,
    battery_remain NUMERIC(5, 2),
    can_be_rent     boolean,
    station_id     int,

    FOREIGN KEY (scooter_id)  REFERENCES scooters (id),
    FOREIGN KEY (location_id) REFERENCES locations (id),
    FOREIGN KEY (station_id)  REFERENCES scooter_stations (id)
);

CREATE TABLE IF NOT EXISTS problem_types
(
    id   smallserial PRIMARY KEY,
    name VARCHAR(150) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS problems
(
    id            bigserial PRIMARY KEY,
    user_id       int       NOT NULL,
    type_Id       smallint  NOT NULL,
    scooter_id    int,
    date_reported TIMESTAMP NOT NULL,
    description   text      NOT NULL,
    is_solved     boolean,

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (type_id) REFERENCES problem_types (id),
    FOREIGN KEY (scooter_id) REFERENCES scooters (id)
);

CREATE TABLE IF NOT EXISTS problem_statuses
(
    problem_id   bigint PRIMARY KEY,
    date_solved  TIMESTAMP NOT NULL,
    description text      NOT NULL,

    FOREIGN KEY (problem_id) REFERENCES problems (id)
);

CREATE TABLE IF NOT EXISTS scooter_statuses_in_rent
(
    id         bigserial PRIMARY KEY,
    user_id     int       NOT NULL,
    scooter_id  int       NOT NULL,
    station_id  int,
    date_time   TIMESTAMP NOT NULL,
    location_id int,
    is_returned boolean,

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (scooter_id) REFERENCES Scooters (id),
    FOREIGN KEY (station_id) REFERENCES Scooter_Stations (id),
    FOREIGN KEY (location_id) REFERENCES Locations (id)
);

CREATE TABLE IF NOT EXISTS orders
(
    id             bigserial PRIMARY KEY,
    user_id        int NOT NULL,
    scooter_id     int NOT NULL,
    status_start_id bigint,
    status_end_id  bigint,
    distance       NUMERIC(12, 2),
    amount         money,

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (scooter_id) REFERENCES scooters (id),
    FOREIGN KEY (status_start_id) REFERENCES scooter_statuses_in_rent (id),
    FOREIGN KEY (status_end_id) REFERENCES scooter_statuses_in_rent (id)
);

CREATE TABLE IF NOT EXISTS account_transactions
(
    id              bigserial PRIMARY KEY,
    date_time       TIMESTAMP NOT NULL,
    payment_type_id smallint  NOT NULL,
    account_from_id int,
    account_to_id   int,
    order_id        bigint,
    amount          money,

    FOREIGN KEY (payment_type_id) REFERENCES payment_types (id),
    FOREIGN KEY (account_from_id) REFERENCES accounts (id),
    FOREIGN KEY (account_To_id) REFERENCES accounts (id),
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

BEGIN;
INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(1, 'admin role', true, false, false);
INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(2, 'user role', false, true, false);
INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(3, 'supplier role', false, false, true);
INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(7, 'super_admin role', true, true, true);

INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('guru_admin@guru.com', false, 'Guru', 'Sadh', 7);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('VikaP@mail.com', false, 'Vika', 'Petrova', 1);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('IraK@mail.com', true, 'Ira', 'Petrova', 1);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('IvanIvanych@mail.com', false, 'Ivan', 'Ivanov', 3);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('PetrPetroff@mail.com', false, 'Petr', 'Petrov', 3);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('UserChan@mail.com', false, 'Jackie', 'Chan', 2);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('UserB@mail.com', true, 'Beyonce', 'Ivanova', 2);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('telo@mail.com', false, 'Goga', 'Boba', 2);
COMMIT;