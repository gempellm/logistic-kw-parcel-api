package internal

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks github.com/gempellm/logistic-kw-parcel-api/internal/app/repo EventRepo
//go:generate mockgen -destination=./mocks/sender_mock.go -package=mocks github.com/gempellm/logistic-kw-parcel-api/internal/app/sender EventSender
