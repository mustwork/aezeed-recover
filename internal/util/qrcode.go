package util

import (
	"github.com/skip2/go-qrcode"
	"strings"
)

func CreateQRCode(content string, padding int) (s string, width, height int, err error) {

	code, err := qrcode.New(content, qrcode.High)
	if err != nil {
		return
	}

	// trimming and then padding seems weird, but we want to keep the border to
	// a minimum and control the padding in order to control the total size.
	bitmap := code.Bitmap()
	bitmap = trimmed(bitmap)
	bitmap = padded(bitmap, padding)
	height = (len(bitmap) / 2) + (len(bitmap) % 2)
	width = len(bitmap[0])

	var builder strings.Builder
	for y := 0; y < len(bitmap)-1; y += 2 {
		for x := 0; x < len(bitmap[y]); x++ {
			top := bitmap[y][x]
			var bottom bool
			if y < len(bitmap) {
				bottom = bitmap[y+1][x]
			} else {
				bottom = false
			}
			if top && bottom {
				builder.WriteRune('\u2588') // full block
			} else if top {
				builder.WriteRune('\u2580') // top half
			} else if bottom {
				builder.WriteRune('\u2584') // bottom half
			} else {
				builder.WriteRune(' ') // empty block
			}
		}
		builder.WriteRune('\n')
	}
	s = builder.String()
	return
}

func padded(bitmap [][]bool, padding int) (result [][]bool) {
	rows := len(bitmap)
	cols := len(bitmap[0])
	result = make([][]bool, rows+padding*2)
	for i := 0; i < padding; i++ {
		result[i] = make([]bool, cols+padding*2)
	}
	for i, row := range bitmap {
		result[i+padding] = make([]bool, cols+padding*2)
		copy(result[i+padding][padding:], row)
	}
	for i := 0; i < padding; i++ {
		result[rows+padding+i] = make([]bool, cols+padding*2)
	}
	return
}

// strips empty rows and columns from the bitmap's edges
func trimmed(bitmap [][]bool) (result [][]bool) {
	// assumes that the bitmap is rectangular
	entireRow := func(row []bool, p bool) bool {
		for _, b := range row {
			if b != p {
				return false
			}
		}
		return true
	}
	entireColumn := func(columnIdx int, m [][]bool, p bool) bool {
		for _, row := range m {
			if row[columnIdx] != p {
				return false
			}
		}
		return true
	}
	rowStart := 0
	rowEnd := len(bitmap)
	for _, row := range bitmap {
		if entireRow(row, false) {
			rowStart++
			continue
		}
		break
	}
	for i := len(bitmap) - 1; i >= 0; i-- {
		if entireRow(bitmap[i], false) {
			rowEnd--
			continue
		}
		break
	}
	colStart := 0
	colEnd := len(bitmap[0])
	for i := 0; i < len(bitmap[0]); i++ {
		if entireColumn(i, bitmap, false) {
			colStart++
			continue
		}
		break
	}
	for i := len(bitmap[0]) - 1; i >= 0; i-- {
		if entireColumn(i, bitmap, false) {
			colEnd--
			continue
		}
		break
	}
	for i := rowStart; i <= rowEnd; i++ {
		result = append(result, bitmap[i][colStart:colEnd])
	}
	return
}
