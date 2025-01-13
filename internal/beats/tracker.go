package beats

type BeatPosition struct {
	Numerator   int
	Denominator int
}

func (bp BeatPosition) Equals(other BeatPosition) bool {
	if bp.Denominator != other.Denominator {
		return false
	}
	return bp.Numerator == other.Numerator
}

type BeatTracker struct {
	Position    BeatPosition
	bpm         float64
	msPerBeat   int64
	currentTime int64
}

func NewBeatTracker(bpm float64, initTime int64) *BeatTracker {
	return &BeatTracker{
		Position: BeatPosition{
			Numerator:   0,
			Denominator: 4,
		},
		bpm:         bpm,
		msPerBeat:   int64((60 * 1000) / bpm),
		currentTime: initTime,
	}
}

func (bt *BeatTracker) Advance(current int64) bool {
	if current < bt.currentTime {
		bt.currentTime = current
	}

	delta := current - bt.currentTime
	if delta >= bt.msPerBeat {
		bt.Position.Numerator = (bt.Position.Numerator + 1) % bt.Position.Denominator
		bt.currentTime = current
		return true
	}
	return false
}

func (bt *BeatTracker) SetCurrentTime(currentTime int64) {
	bt.currentTime = currentTime
}

func (bt *BeatTracker) GetPosition() BeatPosition {
	return bt.Position
}

func (bt *BeatTracker) SetBPM(bpm float64) {
	bt.bpm = bpm
	bt.msPerBeat = int64((60 * 1000) / bpm)
}

func (bt *BeatTracker) SetPosition(pos BeatPosition) {
	bt.Position = pos
}
