module backend

go 1.23.5

require (
	github.com/joho/godotenv v1.5.1
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
)

replace github.com/envoyproxy/go-control-plane/envoy => github.com/envoyproxy/go-control-plane/envoy v1.32.3

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
