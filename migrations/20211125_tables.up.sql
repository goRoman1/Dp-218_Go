CREATE TABLE IF NOT EXISTS Roles
(
    ID         smallint PRIMARY KEY,
    Name       VARCHAR(50) UNIQUE,
    IsAdmin    boolean,
    IsUser     boolean,
    IsSupplier boolean
);

CREATE TABLE IF NOT EXISTS Users
(
    ID          serial PRIMARY KEY,
    LoginEmail  VARCHAR(100) UNIQUE NOT NULL,
    IsBlocked   boolean,
    UserName    VARCHAR(100),
    UserSurname VARCHAR(100),
    CreatedAt   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    RoleID      int                 NOT NULL,

    FOREIGN KEY (RoleID) REFERENCES Roles (ID)
);

CREATE TABLE IF NOT EXISTS LoginInfo
(
    UserId       int PRIMARY KEY,
    PasswordHash VARCHAR(512),

    FOREIGN KEY (UserId) REFERENCES Users (ID)
);

CREATE TABLE IF NOT EXISTS LoginStatus
(
    UserId    int PRIMARY KEY,
    LoggedIn  boolean,
    DateTime  TIMESTAMP NOT NULL,
    IPAddress VARCHAR(40),

    FOREIGN KEY (UserId) REFERENCES Users (ID)
);

CREATE TABLE IF NOT EXISTS ContactTypes
(
    ID   smallserial PRIMARY KEY,
    Name VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS Contacts
(
    ID          serial PRIMARY KEY,
    TypeId      int NOT NULL,
    UserId      int NOT NULL,
    ContactInfo VARCHAR(200),

    FOREIGN KEY (TypeId) REFERENCES ContactTypes (ID),
    FOREIGN KEY (UserId) REFERENCES Users (ID)
);

CREATE TABLE IF NOT EXISTS Accounts
(
    ID      serial PRIMARY KEY,
    Name    VARCHAR(100),
    Number  VARCHAR(100) UNIQUE NOT NULL,
    OwnerId int                 NOT NULL,

    FOREIGN KEY (OwnerId) REFERENCES Users (ID)
);

CREATE TABLE IF NOT EXISTS PaymentTypes
(
    ID   smallserial PRIMARY KEY,
    Name VARCHAR(100) UNIQUE
);

CREATE TABLE IF NOT EXISTS SupplierComissions
(
    ID               serial PRIMARY KEY,
    ComissionPercent NUMERIC(4, 2),
    UserId           int NOT NULL,

    FOREIGN KEY (UserId) REFERENCES Users (ID)
);

CREATE TABLE IF NOT EXISTS SupplierPrices
(
    ID            serial PRIMARY KEY,
    Price         NUMERIC(15, 2),
    PaymentTypeId smallint NOT NULL,
    UserId        int      NOT NULL,

    FOREIGN KEY (PaymentTypeId) REFERENCES PaymentTypes (ID),
    FOREIGN KEY (UserId) REFERENCES Users (ID)
);

CREATE TABLE IF NOT EXISTS ScooterBrands
(
    ID   smallserial PRIMARY KEY,
    Name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS ScooterModels
(
    ID              smallserial PRIMARY KEY,
    BrandId         smallint     NOT NULL,
    PaymentTypeId   smallint     NOT NULL,
    ModelName       VARCHAR(100) NOT NULL,
    BatteryCapacity NUMERIC(10, 0),
    MaxWeight       NUMERIC(5, 2),
    MaxDistance     NUMERIC(10, 0),

    FOREIGN KEY (BrandId) REFERENCES ScooterBrands (ID),
    FOREIGN KEY (PaymentTypeId) REFERENCES PaymentTypes (ID)
);

CREATE TABLE IF NOT EXISTS Scooters
(
    ID           serial PRIMARY KEY,
    ModelId      smallint            NOT NULL,
    OwnerId      int                 NOT NULL,
    SerialNumber VARCHAR(100) UNIQUE NOT NULL,

    FOREIGN KEY (ModelId) REFERENCES ScooterModels (ID),
    FOREIGN KEY (OwnerId) REFERENCES Users (ID)
);

CREATE TABLE IF NOT EXISTS Locations
(
    ID         serial PRIMARY KEY,
    Lattitude  NUMERIC(10, 0) NOT NULL,
    Longtitude NUMERIC(10, 0) NOT NULL,
    Label      VARCHAR(200)
);

CREATE TABLE IF NOT EXISTS ScooterStations
(
    ID         serial PRIMARY KEY,
    LocationId int NOT NULL,
    Name       VARCHAR(100),
    IsActive   boolean,

    FOREIGN KEY (LocationId) REFERENCES Locations (ID)
);

CREATE TABLE IF NOT EXISTS ScooterStatuses
(
    ScooterId     int PRIMARY KEY,
    LocationId    int,
    BatteryRemain NUMERIC(5, 2),
    CanBeRent     boolean,
    StationId     int,

    FOREIGN KEY (ScooterId) REFERENCES Scooters (ID),
    FOREIGN KEY (LocationId) REFERENCES Locations (ID),
    FOREIGN KEY (StationId) REFERENCES ScooterStations (ID)
);

CREATE TABLE IF NOT EXISTS ProblemTypes
(
    ID   smallserial PRIMARY KEY,
    Name VARCHAR(150) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS Problems
(
    ID           bigserial PRIMARY KEY,
    UserId       int       NOT NULL,
    TypeId       smallint  NOT NULL,
    ScooterId    int,
    DateReported TIMESTAMP NOT NULL,
    Description  text      NOT NULL,
    IsSolved     boolean,

    FOREIGN KEY (UserId) REFERENCES Users (ID),
    FOREIGN KEY (TypeId) REFERENCES ProblemTypes (ID),
    FOREIGN KEY (ScooterId) REFERENCES Scooters (ID)
);

CREATE TABLE IF NOT EXISTS ProblemStatuses
(
    ProblemID   bigint PRIMARY KEY,
    DateSolved  TIMESTAMP NOT NULL,
    Description text      NOT NULL,

    FOREIGN KEY (ProblemID) REFERENCES Problems (ID)
);

CREATE TABLE IF NOT EXISTS ScooterStatusesInRent
(
    ID         bigserial PRIMARY KEY,
    UserId     int       NOT NULL,
    ScooterId  int       NOT NULL,
    StationId  int,
    DateTime   TIMESTAMP NOT NULL,
    LocationId int,
    IsReturned boolean,

    FOREIGN KEY (UserId) REFERENCES Users (ID),
    FOREIGN KEY (ScooterId) REFERENCES Scooters (ID),
    FOREIGN KEY (StationId) REFERENCES ScooterStations (ID),
    FOREIGN KEY (LocationId) REFERENCES Locations (ID)
);

CREATE TABLE IF NOT EXISTS Orders
(
    ID            bigserial PRIMARY KEY,
    UserId        int NOT NULL,
    ScooterId     int NOT NULL,
    StatusStartId bigint,
    StatusEndId   bigint,
    Distance      NUMERIC(12, 2),
    Amount        money,

    FOREIGN KEY (UserId) REFERENCES Users (ID),
    FOREIGN KEY (ScooterId) REFERENCES Scooters (ID),
    FOREIGN KEY (StatusStartId) REFERENCES ScooterStatusesInRent (ID),
    FOREIGN KEY (StatusEndId) REFERENCES ScooterStatusesInRent (ID)
);

CREATE TABLE IF NOT EXISTS AccountTransactions
(
    ID            bigserial PRIMARY KEY,
    DateTime      TIMESTAMP NOT NULL,
    PaymentTypeId smallint  NOT NULL,
    AccountFromId int,
    AccountToId   int,
    OrderId       bigint,
    Amount        money,

    FOREIGN KEY (PaymentTypeId) REFERENCES PaymentTypes (ID),
    FOREIGN KEY (AccountFromId) REFERENCES Accounts (ID),
    FOREIGN KEY (AccountToId) REFERENCES Accounts (ID),
    FOREIGN KEY (OrderId) REFERENCES Orders (ID)
);

BEGIN;
INSERT INTO Roles(ID, Name, IsAdmin, IsUser, IsSupplier) VALUES(1, 'Admin role', true, false, false);
INSERT INTO Roles(ID, Name, IsAdmin, IsUser, IsSupplier) VALUES(2, 'User role', false, true, false);
INSERT INTO Roles(ID, Name, IsAdmin, IsUser, IsSupplier) VALUES(3, 'Supplier role', false, false, true);
INSERT INTO Roles(ID, Name, IsAdmin, IsUser, IsSupplier) VALUES(7, 'SuperAdmin role', true, true, true);
COMMIT;