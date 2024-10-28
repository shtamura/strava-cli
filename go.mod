module github.com/shtamura/strava-cli

go 1.22.2

require github.com/spf13/cobra v1.8.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/shtamura/oauth2 v0.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/shtamura/oauth2 => ./pkg/oauth2
