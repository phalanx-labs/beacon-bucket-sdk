package utility

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TimestampToTime 将 protobuf Timestamp 转换为 time.Time
func TimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// OptionalTimestampToTime 将可选的 protobuf Timestamp 转换为 *time.Time
func OptionalTimestampToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}
