package validations

import "errors"

var (
	ErrInvalidServiceName    = errors.New("service name must be provided")
	ErrInvalidSubscriptionID = errors.New("invalid subscription ID")
	ErrInvalidPrice          = errors.New("price must be positive integer")
	ErrInvalidDateFormat     = errors.New("invalid date format, expected MM-YYYY")
	ErrEndDateBeforeStart    = errors.New("end date must not be lessthan start date")
	ErrInvalidUserID         = errors.New("invalid user ID")
	ErrEmptyUserID           = errors.New("user ID is empty")
	ErrSubscriptionExists    = errors.New("subscription already exists")
	ErrSubscriptionNotFound  = errors.New("subscription not found")
	ErrInvalidStartDate      = errors.New("invalid start_date format, expected MM-YYYY")
	ErrInvalidEndDate        = errors.New("invalid end_date format, expected MM-YYYY")
	ErrInvalidRequestInput   = errors.New("invalid request input")
	ErrInvalid               = errors.New("invalid query parameters")
	//Repo Error
	ErrCreateSubscriptionFailed       = errors.New("failed to create subscription")
	ErrListSubscriptionFailed         = errors.New("failed to list subscription")
	ErrGetSubscriptionByIDFailed      = errors.New("failed to get subscription by ID")
	ErrUpdateSubscriptionFailed       = errors.New("failed to get update subscription")
	ErrDeleteSubscriptionFailed       = errors.New("failed to delete subscription")
	ErrCalculateTotalCostFailed       = errors.New("failed to calculate totalcost")
	ErrFindSubscriptionByPeriodFailed = errors.New("failed to find subscription by userId or servicename")
	//Database Error
	ErrDbInitializationFailed  = errors.New("failed to initialize db")
	ErrDbMigrationFailed       = errors.New("migration failed")
	ErrDbConnectionFailed      = errors.New("failed to connect to database")
	ErrDbPingFailed            = errors.New("failed to ping db")
	ErrDbCloseConnectionFailed = errors.New("failed to close database connections")
	//Config Error
	ErrConfiLoadFailed = errors.New("failed to load config from environment, config set to default value")

	//router error
	ErrServerStartFailed = errors.New("failed to start the server.")
	//AppErrr
	ErrInvalidGinMode       = errors.New("Invalid GIN_MODE")
	ErrShuttingServerFailed = errors.New("error during server shutdown.")
)
