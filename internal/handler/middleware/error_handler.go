package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.socket9.com/sritrang/super-driver/backend/internal/core/domain"
	"gitlab.socket9.com/sritrang/super-driver/backend/internal/handler/common"
)

// ErrorHandler returns an Echo HTTPErrorHandler that maps domain and handler errors to JSON.
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		reqID := c.Response().Header().Get(echo.HeaderXRequestID)
		var httpStatus int
		var resp common.ErrorResponse

		switch e := err.(type) {
		case *common.Error:
			httpStatus = e.HTTPStatus
			resp = common.ErrorResponse{
				IsError:   true,
				Code:      e.Code,
				Message:   e.Message,
				RequestID: reqID,
				ErrorType: e.ErrorType,
				Errors:    e.Errors,
			}
		case *domain.Error:
			httpStatus = domainHTTPStatus(e.Code)
			resp = common.ErrorResponse{
				IsError:   true,
				Code:      common.GetAppPrefix() + "-" + string(e.Code),
				Message:   e.Message,
				RequestID: reqID,
				ErrorType: "DOMAIN_ERROR",
				Errors:    e.Errors,
			}
		case *echo.HTTPError:
			httpStatus = e.Code
			resp = common.ErrorResponse{
				IsError:   true,
				Code:      common.GetAppPrefix() + "-" + echoCode(e.Code),
				Message:   http.StatusText(e.Code),
				RequestID: reqID,
			}
		default:
			httpStatus = http.StatusInternalServerError
			resp = common.ErrorResponse{
				IsError:   true,
				Code:      common.GetAppPrefix() + "-5001",
				Message:   "Internal Server Error",
				RequestID: reqID,
			}
		}

		_ = c.JSON(httpStatus, resp)
	}
}

func domainHTTPStatus(code domain.ErrorCode) int {
	switch code {
	case domain.ErrorCodeBadRequest:
		return http.StatusBadRequest
	case domain.ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrorCodeForbidden:
		return http.StatusForbidden
	case domain.ErrorCodeNotFound:
		return http.StatusNotFound
	case domain.ErrorCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func echoCode(status int) string {
	switch status {
	case 400:
		return "4001"
	case 401:
		return "4011"
	case 403:
		return "4031"
	case 404:
		return "4041"
	case 409:
		return "4091"
	default:
		return "5001"
	}
}
