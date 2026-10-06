package domain

type ArrowScore struct {
	Code         string
	NumericValue int
}

func ParseArrowScore(code string) (ArrowScore, error) {
	value, ok := numericValue(code)
	if !ok {
		return ArrowScore{}, ErrInvalidScoreCode
	}
	return ArrowScore{Code: code, NumericValue: value}, nil
}

func (s ArrowScore) IsGold() bool {
	switch s.Code {
	case "X", "10", "9":
		return true
	case "8", "7", "6", "5", "4", "3", "2", "1", "M":
		return false
	default:
		panic("unparsed score code: " + s.Code)
	}
}

func (s ArrowScore) IsHit() bool {
	return s.Code != "M"
}

func numericValue(code string) (int, bool) {
	switch code {
	case "X":
		return 10, true
	case "10":
		return 10, true
	case "9":
		return 9, true
	case "8":
		return 8, true
	case "7":
		return 7, true
	case "6":
		return 6, true
	case "5":
		return 5, true
	case "4":
		return 4, true
	case "3":
		return 3, true
	case "2":
		return 2, true
	case "1":
		return 1, true
	case "M":
		return 0, true
	default:
		return 0, false
	}
}
