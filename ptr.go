package invezgo

// Ptr returns a pointer to v; use it for optional params fields:
//
//	&invezgo.CalendarParams{Page: invezgo.Ptr(int64(2))}
func Ptr[T any](v T) *T { return &v }
