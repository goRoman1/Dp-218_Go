package configs

import (
	"os"
	"strconv"
)

var PG_HOST = os.Getenv("PG_HOST")
var PG_PORT = os.Getenv("PG_PORT")
var POSTGRES_DB = os.Getenv("POSTGRES_DB")
var POSTGRES_USER = os.Getenv("POSTGRES_USER")
var POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
var HTTP_PORT = os.Getenv("HTTP_PORT")
var MIGRATE_DOWN, _ = strconv.ParseBool(os.Getenv("MIGRATE_DOWN"))
var MIGRATIONS_PATH = os.Getenv("MIGRATIONS_PATH")
var TEMPLATES_PATH = os.Getenv("TEMPLATES_PATH")
var MIGRATE_VERSION_FORCE, _ = strconv.Atoi(os.Getenv("MIGRATE_VERSION_FORCE"))
var KAFKA_BROKER = os.Getenv("KAFKA_BROKER")
var SESSION_SECRET = os.Getenv("SESSION_SECRET")
