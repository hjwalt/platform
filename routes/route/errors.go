package route

import "errors"

var (
	Err400 = errors.New("bad request")
	Err401 = errors.New("unauthorised")
	Err404 = errors.New("not found")
	Err500 = errors.New("server errror")
)

func ErrorCode(err error) int {
	if err == nil {
		return 0
	}

	if errors.Is(err, Err400) {
		return 400
	}
	if errors.Is(err, Err401) {
		return 401
	}
	if errors.Is(err, Err404) {
		return 404
	}
	if errors.Is(err, Err500) {
		return 500
	}

	return 500
}
