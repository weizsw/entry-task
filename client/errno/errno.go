package errno

const (
	StatusOK          = 0
	StatusParamsError = 400
	StatusServerError = 500

	LoginCredentialError = 11000
	RegisterFailedError  = 11001
)

var ErrorMsg = map[int]string{
	StatusOK:             "ok",
	StatusParamsError:    "params error",
	StatusServerError:    "server error",
	LoginCredentialError: "login credentials error",
	RegisterFailedError:  "something went wrong",
}
