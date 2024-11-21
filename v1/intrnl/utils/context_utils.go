package utils

import "github.com/labstack/echo/v4"

func GetValueFromContext[T any](c echo.Context, key string) (T, bool) {
	value := c.Get(key)
	if value == nil {
		var zeroVal T
		return zeroVal, false
	}

	v, ok := value.(T)

	return v, ok
}

func SetValueInContext[T any](c echo.Context, key string, value T) {
	c.Set(key, value)
}
