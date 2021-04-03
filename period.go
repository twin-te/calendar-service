package main

func GetPeriodStart(period int) (int, int) {
	switch period {
	case 1:
		return 8, 40
	case 2:
		return 10, 10
	case 3:
		return 12, 15
	case 4:
		return 13, 45
	case 5:
		return 15, 15
	case 6:
		return 16, 45
	}
	return 0, 0
}

func GetPeriodEnd(period int) (int, int) {
	switch period {
	case 1:
		return 9, 55
	case 2:
		return 11, 25
	case 3:
		return 13, 30
	case 4:
		return 15, 00
	case 5:
		return 16, 30
	case 6:
		return 18, 00
	}
	return 23, 59
}
