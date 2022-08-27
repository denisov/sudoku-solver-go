package main

// PossibleValues это возможные значения в данной ячейке
type PossibleValues struct {
	column int
	row    int
	values []int
}

func solve(square *Square) bool {
	var cellWithMinPossibleValues *PossibleValues
	for {
		cellWithMinPossibleValues = getCellWithMinPossibleValues(*square)
		if cellWithMinPossibleValues == nil {
			return true
		}
		if len((*cellWithMinPossibleValues).values) == 1 {
			column := (*cellWithMinPossibleValues).column
			row := (*cellWithMinPossibleValues).row
			square[row][column] = (*cellWithMinPossibleValues).values[0]
			continue
		}
		if len((*cellWithMinPossibleValues).values) == 0 {
			return false
		}
		if len((*cellWithMinPossibleValues).values) > 1 {
			break
		}
	}

	results := make(chan *Square, len((*cellWithMinPossibleValues).values))
	for _, v := range (*cellWithMinPossibleValues).values {
		go func(value int, res chan<- *Square) {
			square_copy := *square
			square_copy[cellWithMinPossibleValues.row][cellWithMinPossibleValues.column] = value
			if solve(&square_copy) {
				res <- &square_copy
			} else {
				res <- nil
			}
		}(v, results)

	}

	for n := 0; n < len((*cellWithMinPossibleValues).values); n++ {
		squareCandidate := <-results
		if squareCandidate != nil {
			*square = *squareCandidate
			return true
		}
	}

	return false
}

// getCellWithMinPossibleValues возвращает ячейку с наименьшим числом возможных значений
func getCellWithMinPossibleValues(square Square) *PossibleValues {
	var possibleValues *PossibleValues
	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			if square[row][column] != 0 {
				continue
			}
			possibleValuesCandidate := getPossibleValues(column, row, square)
			if possibleValues == nil || len(possibleValuesCandidate) < len((*possibleValues).values) {
				possibleValues = &PossibleValues{
					values: possibleValuesCandidate,
					column: column,
					row:    row,
				}
			}
		}
	}

	return possibleValues
}

// getPossibleValues Возвращает возможные значения для данной ячейки
func getPossibleValues(column int, row int, square Square) []int {

	existingValues := make(map[int]bool, 9)
	for i := 1; i < 10; i++ {
		existingValues[i] = false
	}

	// строка
	for _, v := range square[row] {
		existingValues[v] = true
	}

	// столбец
	for r := 0; r < 9; r++ {
		val := square[r][column]
		existingValues[val] = true
	}

	for _, v := range getValuesFromSubsquare(square, (column/3)*3, (row/3)*3) {
		existingValues[v] = true
	}

	possibleValues := []int{}
	for key, val := range existingValues {
		if !val {
			possibleValues = append(possibleValues, key)
		}
	}

	return possibleValues
}

// getValuesFromSubsquare Возвращает значения из под-квадрата
func getValuesFromSubsquare(square Square, columnStart, rowStart int) []int {
	size := 3
	var values []int
	for row := rowStart; row < rowStart+size; row++ {
		for column := columnStart; column < columnStart+size; column++ {
			values = append(values, square[row][column])
		}
	}
	return values
}
