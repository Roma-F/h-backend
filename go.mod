module lbr-backend

go 1.18

require (
	github.com/blockloop/scan v1.3.0
	github.com/codingsince1985/geo-golang v1.8.2
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/qustavo/dotsql v1.1.0
	github.com/rs/zerolog v1.28.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/doug-martin/goqu/v9 v9.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	golang.org/x/sys v0.0.0-20221013171732-95e765b1cc43 // indirect
)

replace github.com/codingsince1985/geo-golang v1.8.2 => github.com/NetScrn/geo-golang v0.0.0-20221119122320-01d7157f0ab8
