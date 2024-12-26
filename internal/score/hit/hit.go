package hit

// Rating is the score rating of a hit
type HitRating string

type hitRating struct {
	Perfect HitRating
	Good    HitRating
	Bad     HitRating
	Miss    HitRating
}

// Window is the number of ms in which a note can be hit
type hitWindow struct {
	Perfect int64
	Good    int64
	Bad     int64
}

type hitValue struct {
	Perfect int
	Good    int
	Bad     int
	Miss    int
}

var (
	Rating = hitRating{
		Perfect: "Perfect",
		Good:    "Good",
		Bad:     "Bad",
		Miss:    "Miss",
	}

	Value = hitValue{
		Perfect: 10,
		Good:    5,
		Bad:     2,
		Miss:    0,
	}

	Window = hitWindow{
		Perfect: 40,
		Good:    50,
		Bad:     60,
	}
)

func GetHitRating(diff int64) HitRating {
	if diff < Window.Perfect {
		return Rating.Perfect
	} else if diff < Window.Good {
		return Rating.Good
	} else if diff < Window.Bad {
		return Rating.Bad
	}
	return Rating.Miss
}
