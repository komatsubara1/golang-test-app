package request

type DebugSetTimeShiftRequest struct {
	TimeShift int64 `json:"time_shift" binding:"required"`
}
