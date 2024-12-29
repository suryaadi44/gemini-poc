// File: ptr.go
//
// Defines custom helper type such as primitive to pointer conversion.
// This is useful for when we want to pass a primitive type as a pointer
// to a function that requires a pointer to function as intended such as
// gorm's Updates() function that requires a pointer to a struct to update
// to default values.
//
// Surya Adi (kmsurya.adi44@gmail.com)

package custom

func ToPtr[T any](v T) *T {
	return &v
}

func ToVal[T any](v *T) T {
	return *v
}

func ToValOrDefault[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}
