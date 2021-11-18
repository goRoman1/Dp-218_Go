DROP TABLE IF EXISTS Roles CASCADE;
CREATE TABLE Roles(
    ID      serial PRIMARY KEY,
    Name    VARCHAR(50) UNIQUE,
    IsAdmin boolean,
    IsUser  boolean,
    IsSupplier boolean
);

DROP TABLE IF EXISTS Users CASCADE;
CREATE TABLE Users(
    ID          serial PRIMARY KEY,
    LoginEmail  VARCHAR(100) UNIQUE NOT NULL,
    IsBlocked   boolean,
    UserName    VARCHAR(100),
    UserSurname VARCHAR(100),
    CreatedAt   TIMESTAMP NOT NULL,
    RoleID      int NOT NULL,

    FOREIGN KEY (RoleID) REFERENCES Roles(ID)
);


DROP TABLE IF EXISTS LoginInfo CASCADE;
CREATE TABLE LoginInfo(
    UserId  int  PRIMARY KEY,
    PasswordHash VARCHAR(512),

    FOREIGN KEY (UserId) REFERENCES Users(ID)
);

DROP TABLE IF EXISTS LoginStatus CASCADE;
CREATE TABLE LoginStatus(
    UserId  int  PRIMARY KEY,
    LoggedIn boolean,
    DateTime TIMESTAMP NOT NULL,
    IPAddress VARCHAR(40),

    FOREIGN KEY (UserId) REFERENCES Users(ID)
);

DROP TABLE IF EXISTS ContactTypes CASCADE;
CREATE TABLE ContactTypes(
    ID  smallserial PRIMARY KEY,
    Name VARCHAR(50)
);

DROP TABLE IF EXISTS Contacts CASCADE;
CREATE TABLE Contacts(
    ID  serial PRIMARY KEY,
    TypeId int NOT NULL,
    UserId int NOT NULL,
    ContactInfo VARCHAR(200),

    FOREIGN KEY (TypeId) REFERENCES ContactTypes(ID),
    FOREIGN KEY (UserId) REFERENCES Users(ID)
);

DROP TABLE IF EXISTS Accounts CASCADE;
CREATE TABLE Accounts(
    ID serial PRIMARY KEY,
    Name VARCHAR(100),
    Number VARCHAR(100) UNIQUE NOT NULL,
    OwnerId int NOT NULL,

    FOREIGN KEY (OwnerId) REFERENCES Users(ID)
);

DROP TABLE IF EXISTS PaymentTypes CASCADE;
CREATE TABLE PaymentTypes(
    ID smallserial PRIMARY KEY,
    Name VARCHAR(100) UNIQUE
);

DROP TABLE IF EXISTS SupplierComissions CASCADE;
CREATE TABLE SupplierComissions(
    ID serial PRIMARY KEY,
    ComissionPercent NUMERIC(4,2),
    UserId int NOT NULL,

    FOREIGN KEY (UserId) REFERENCES Users(ID)
);

DROP TABLE IF EXISTS SupplierPrices CASCADE;
CREATE TABLE SupplierPrices(
    ID serial PRIMARY KEY,
    Price NUMERIC(15,2),
    PaymentTypeId smallint NOT NULL,
    UserId int NOT NULL,

    FOREIGN KEY (PaymentTypeId) REFERENCES PaymentTypes(ID),
    FOREIGN KEY (UserId) REFERENCES Users(ID)
);

DROP TABLE IF EXISTS ScooterBrands CASCADE;
CREATE TABLE ScooterBrands(
    ID smallserial PRIMARY KEY,
    Name VARCHAR(100) UNIQUE NOT NULL
);

DROP TABLE IF EXISTS ScooterModels CASCADE;
CREATE TABLE ScooterModels(
    ID smallserial PRIMARY KEY,
    BrandId smallint NOT NULL,
    PaymentTypeId smallint NOT NULL,
    ModelName VARCHAR(100) NOT NULL,
    BatteryCapacity NUMERIC(10, 0),
    MaxWeight NUMERIC(5, 2),
    MaxDistance NUMERIC(10, 0),

    FOREIGN KEY (BrandId) REFERENCES ScooterBrands(ID),
    FOREIGN KEY (PaymentTypeId) REFERENCES PaymentTypes(ID)
);

DROP TABLE IF EXISTS Scooters CASCADE;
CREATE TABLE Scooters(
    ID serial PRIMARY KEY,
    ModelId smallint NOT NULL,
    OwnerId int NOT NULL,
    SerialNumber VARCHAR(100) UNIQUE NOT NULL,

    FOREIGN KEY (ModelId) REFERENCES ScooterModels(ID),
    FOREIGN KEY (OwnerId) REFERENCES Users(ID)
); 

DROP TABLE IF EXISTS Locations CASCADE;
CREATE TABLE Locations(
    ID serial PRIMARY KEY,
    Lattitude NUMERIC(10,0) NOT NULL,
    Longtitude NUMERIC(10,0) NOT NULL,
    Label VARCHAR(200)
);

DROP TABLE IF EXISTS ScooterStations CASCADE;
CREATE TABLE ScooterStations(
    ID serial PRIMARY KEY,
    LocationId int NOT NULL,
    Name VARCHAR(100),
    IsActive boolean,

    FOREIGN KEY (LocationId) REFERENCES Locations(ID)
);

DROP TABLE IF EXISTS ScooterStatuses CASCADE;
CREATE TABLE ScooterStatuses(
    ScooterId int PRIMARY KEY,
    LocationId int,
    BatteryRemain NUMERIC(5,2),
    CanBeRent boolean,
    StationId int,
    
    FOREIGN KEY (ScooterId) REFERENCES Scooters(ID),
    FOREIGN KEY (LocationId) REFERENCES Locations(ID),
    FOREIGN KEY (StationId) REFERENCES ScooterStations(ID)
); 

DROP TABLE IF EXISTS ProblemTypes CASCADE;
CREATE TABLE ProblemTypes(
    ID smallserial PRIMARY KEY,
    Name VARCHAR(150) UNIQUE NOT NULL
);

DROP TABLE IF EXISTS Problems CASCADE;
CREATE TABLE Problems(
    ID bigserial PRIMARY KEY,
    UserId int NOT NULL,
    TypeId smallint NOT NULL,
    ScooterId int,
    DateReported TIMESTAMP NOT NULL,
    Description text NOT NULL,
    IsSolved boolean,

    FOREIGN KEY (UserId) REFERENCES Users(ID),
    FOREIGN KEY (TypeId) REFERENCES ProblemTypes(ID),
    FOREIGN KEY (ScooterId) REFERENCES Scooters(ID)
);

DROP TABLE IF EXISTS Problems CASCADE;
CREATE TABLE Problems(
    ProblemID bigint PRIMARY KEY,
    DateSolved TIMESTAMP NOT NULL,
    Description text NOT NULL,

    FOREIGN KEY (ProblemID) REFERENCES Problems(ID)
);

DROP TABLE IF EXISTS ScooterStatusesInRent CASCADE;
CREATE TABLE ScooterStatusesInRent(
    ID bigserial PRIMARY KEY,
    UserId int NOT NULL,
    ScooterId int NOT NULL,
    StationId int,
    DateTime TIMESTAMP NOT NULL,
    LocationId int,
    IsReturned boolean,

    FOREIGN KEY (UserId) REFERENCES Users(ID),
    FOREIGN KEY (ScooterId) REFERENCES Scooters(ID),
    FOREIGN KEY (StationId) REFERENCES ScooterStations(ID),
    FOREIGN KEY (LocationId) REFERENCES Locations(ID)
);

DROP TABLE IF EXISTS Orders CASCADE;
CREATE TABLE Orders(
    ID bigserial PRIMARY KEY,
    UserId int NOT NULL,
    ScooterId int NOT NULL,
    StatusStartId bigint,
    StatusEndId bigint,
    Distance NUMERIC(12,2),
    Amount money,

    FOREIGN KEY (UserId) REFERENCES Users(ID),
    FOREIGN KEY (ScooterId) REFERENCES Scooters(ID),
    FOREIGN KEY (StatusStartId) REFERENCES ScooterStatusesInRent(ID),
    FOREIGN KEY (StatusEndId) REFERENCES ScooterStatusesInRent(ID)
);

DROP TABLE IF EXISTS AccountTransactions CASCADE;
CREATE TABLE AccountTransactions(
    ID bigserial PRIMARY KEY,
    DateTime TIMESTAMP NOT NULL,
    PaymentTypeId smallint NOT NULL,
    AccountFromId int,
    AccountToId int,
    OrderId bigint,
    Amount money,

    FOREIGN KEY (PaymentTypeId) REFERENCES PaymentTypes(ID),
    FOREIGN KEY (AccountFromId) REFERENCES Accounts(ID),
    FOREIGN KEY (AccountToId) REFERENCES Accounts(ID),
    FOREIGN KEY (OrderId) REFERENCES Orders(ID)
);