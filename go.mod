module github.com/gempellm/logistic-parcel-api

go 1.16

replace github.com/gempellm/logistic-parcel-api/pkg/logistic_parcel_api => ./pkg/logistic_parcel_api

replace github.com/gempellm/logistic-kw-parcel-api/internal/app/retranslator => ./internal/app/retranslator

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/gammazero/workerpool v1.1.3
	github.com/gempellm/logistic-kw-parcel-api v0.0.0-20220928153703-ffe6a2e42b0b
	github.com/gempellm/logistic-parcel-api/pkg/logistic_parcel_api v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.12.0
	github.com/jackc/pgx/v4 v4.17.2
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.7
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pressly/goose/v3 v3.7.0
	github.com/prometheus/client_golang v1.13.0
	github.com/rs/zerolog v1.28.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	google.golang.org/grpc v1.50.1
	gopkg.in/yaml.v3 v3.0.1
)
