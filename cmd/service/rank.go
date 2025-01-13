package main

const MaxScore = 100000

func getRatingValue(s *Score) float64 {
	scorePerc := float64(s.Score) / float64(100000)
	return scorePerc * float64(s.Difficulty)
}
