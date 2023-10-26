package main

func sliceFill[T any](arr *[]T, value T) {
	for i := range *arr {
		(*arr)[i] = value
	}
}
func mt_sliceFill[T any](arr *[][]T, value T) {
	for x := range *arr {
		for y := range (*arr)[x] {
			(*arr)[x][y] = value
		}
	}
}

// faulty!!!
// func slice_make[T any](size int, initial_value T) []T {
// 	slice := make([]T, size)
// 	for i := range slice {
// 		slice[i] = initial_value
// 	}
// 	return slice
// }

func init_mdimension[T any](slice *[][]T, size int) {
	for w := range *slice {
		(*slice)[w] = make([]T, size)
	}
}
