package main

func printCode(i int) (string) {
	switch i {
	case 0:
		return "_"
	case 1:
		return "X"
	case 2:
		return "0"
	default:
		return "_"
	}
}

func getTable(room *Room) (string, string, string) {
	var line1 string = " "
	var line2 string = " "
	var line3 string = " "
	for i, val := range room.cells {
		switch i {
		case 0, 1, 2:
			line1 += printCode(val) + " "
		case 3, 4, 5:
			line2 += printCode(val) + " "
		case 6, 7, 8:
			line3 += printCode(val) + " "
		}
	}
	return line1, line2, line3
}

func findWinner(room *Room) (int) {
	//Horizontal
	if room.cells[0] == room.cells[1] && room.cells[1] == room.cells[2] && room.cells[0] != 0 {
		return room.cells[0]
	}
	if room.cells[3] == room.cells[4] && room.cells[4] == room.cells[5] && room.cells[3] != 0 {
		return room.cells[3]
	}
	if room.cells[6] == room.cells[7] && room.cells[7] == room.cells[8] && room.cells[6] != 0 {
		return room.cells[6]
	}
	//Vertical
	if room.cells[0] == room.cells[3] && room.cells[3] == room.cells[6] && room.cells[0] != 0 {
		return room.cells[0]
	}
	if room.cells[1] == room.cells[4] && room.cells[4] == room.cells[7] && room.cells[1] != 0 {
		return room.cells[1]
	}
	if room.cells[2] == room.cells[5] && room.cells[5] == room.cells[8] && room.cells[2] != 0 {
		return room.cells[2]
	}
	//Diagonal
	if room.cells[0] == room.cells[4] && room.cells[4] == room.cells[8] && room.cells[0] != 0 {
		return room.cells[0]
	}
	if room.cells[2] == room.cells[4] && room.cells[4] == room.cells[6] && room.cells[2] != 0 {
		return room.cells[2]
	}
	return 0
}
