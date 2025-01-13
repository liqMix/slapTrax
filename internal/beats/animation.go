package beats

type PulseAnimation struct {
	scale       float64
	targetScale float64
	speed       float64
}

func NewPulseAnimation(targetScale float64, speed float64) *PulseAnimation {
	return &PulseAnimation{
		scale:       1.0,
		targetScale: targetScale,
		speed:       speed,
	}
}

func (la *PulseAnimation) Pulse() {
	la.scale = la.targetScale
}

func (la *PulseAnimation) Update() {
	if la.scale > 1.0 {
		la.scale -= la.speed
		if la.scale < 1.0 {
			la.scale = 1.0
		}
	}
}

func (la *PulseAnimation) Reset() {
	la.scale = 1.0
}

func (la *PulseAnimation) GetScale() float64 {
	return la.scale
}
