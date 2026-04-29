package utils

import "strconv"

func Ptr[T any](v T) *T {
    return &v
}

func ToInt(valStr string, defaultVal int) int {
    if valStr == "" {
        return defaultVal
    }

    val, err := strconv.Atoi(valStr)
    if err != nil {
        return -2
    }
    return val
}