package types

import (
	"image/color"
)

type Rank string

const (
	SlapNewbie     Rank = "slapChump"
	SlapRookie     Rank = "slapBum"
	SlapContender  Rank = "slapWannabe"
	SlapBrawler    Rank = "slapHappy"
	SlapChampion   Rank = "slapChamp"
	SlapElite      Rank = "slapElite"
	SlapLegend     Rank = "slapLegend"
	SlapUndisputed Rank = "slapPope"
	SlapGod        Rank = "slapGod"
)

func (r Rank) String() string {
	return string(r)
}

func (r Rank) Color() color.RGBA {
	switch r {
	case SlapNewbie:
		return color.RGBA{255, 255, 255, 255} // #ffffff
	case SlapRookie:
		return color.RGBA{30, 175, 30, 255} // #00ff00
	case SlapContender:
		return color.RGBA{0, 0, 255, 255} // #0000ff
	case SlapBrawler:
		return color.RGBA{255, 0, 255, 255} // #ff00ff
	case SlapChampion:
		return color.RGBA{255, 255, 0, 255} // #ffff00
	case SlapElite:
		return color.RGBA{0, 255, 255, 255} // #00ffff
	case SlapLegend:
		return color.RGBA{255, 128, 0, 255} // #ff8000
	case SlapUndisputed:
		return color.RGBA{128, 0, 255, 255} // #8000ff
	case SlapGod:
		return color.RGBA{255, 0, 0, 255} // #ff0000
	}
	return color.RGBA{255, 255, 255, 255} // #ffffff
}

func RankTitleFromRank(rank float64) Rank {
	if rank < 5 {
		return SlapNewbie
	} else if rank < 10 {
		return SlapRookie
	} else if rank < 15 {
		return SlapContender
	} else if rank < 20 {
		return SlapBrawler
	} else if rank < 30 {
		return SlapChampion
	} else if rank < 40 {
		return SlapElite
	} else if rank < 50 {
		return SlapLegend
	} else if rank < 60 {
		return SlapUndisputed
	} else {
		return SlapGod
	}
}
