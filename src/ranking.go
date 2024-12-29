package src

type RANK = int64

const RANK_0 RANK = 0
const RANK_1 RANK = 1
const RANK_2 RANK = 2
const RANK_3 RANK = 3
const RANK_4 RANK = 4
const RANK_5 RANK = 5
const RANK_6 RANK = 6
const RANK_7 RANK = 7
const RANK_8 RANK = 8
const RANK_9 RANK = 9
const RANK_PICTURE RANK = -1

type HALFPAYOUTMODE = RANK

const RANK_NONE HALFPAYOUTMODE = 10

func Points(hand CardGetter) int64 {
	var sum int64
	for _, c := range hand.GetCards(3) {
		sum += c.FaceValue
	}
	sum %= 10
	return sum
}

// CheckUpperHand
// returns nil values on draw
func CheckUpperHand(hand1 *Hand, hand2 *Hand) (winner *Hand, loser *Hand) {
	if IsTriplePicture(hand1) && IsTriplePicture(hand2) {
		return nil, nil
	}
	if IsTriplePicture(hand1) {
		return hand1, hand2
	}
	if IsTriplePicture(hand2) {
		return hand2, hand1
	}

	if Points(hand1) > Points(hand2) {
		return hand1, hand2
	}
	if Points(hand1) < Points(hand2) {
		return hand2, hand1
	}

	if isDoublePicture(hand1) && isDoublePicture(hand2) {
		return nil, nil
	}
	if isDoublePicture(hand1) {
		return hand1, hand2
	}
	if isDoublePicture(hand2) {
		return hand2, hand1
	}
	if isSinglePicture(hand1) && isSinglePicture(hand2) {
		return nil, nil
	}
	if isSinglePicture(hand1) {
		return hand1, hand2
	}
	if isSinglePicture(hand2) {
		return hand2, hand1
	}

	return nil, nil
}

func IsTriplePicture(hand CardGetter) bool {
	for _, c := range hand.GetCards(3) {
		if !c.Picture {
			return false
		}
	}
	return true
}

func isDoublePicture(hand CardGetter) bool {
	t := 0
	for _, c := range hand.GetCards(3) {
		if c.Picture {
			t += 1
		}
	}
	return t == 2
}

func isSinglePicture(hand CardGetter) bool {
	t := 0
	for _, c := range hand.GetCards(3) {
		if c.Picture {
			t += 1
		}
	}
	return t == 1
}
