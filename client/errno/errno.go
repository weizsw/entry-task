package errno

const (
	StatusOK          = 0
	StatusParamsError = 400
	StatusAuthError   = 401
	StatusServerError = 500

	LoginCredentialError = 11000
	RegisterFailedError  = 11001
	RedisError           = 12000
	MySQLError           = 13000
)

var ErrorMsg = map[int]string{
	StatusOK:             "ok",
	StatusParamsError:    "params error",
	StatusAuthError:      "auth error",
	StatusServerError:    "server error",
	LoginCredentialError: "login credentials error",
	RegisterFailedError:  "something went wrong",
	RedisError:           "redis error",
	MySQLError:           "mysql error",
}
