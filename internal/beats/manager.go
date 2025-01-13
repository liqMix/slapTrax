package beats

type Manager struct {
	tracker    *BeatTracker
	triggerMap *BeatTriggerMap
}

func NewManager(bpm float64, initMs int64) *Manager {
	return &Manager{
		tracker:    NewBeatTracker(bpm, initMs),
		triggerMap: NewBeatTriggerMap(),
	}
}

func (m *Manager) Update(currentTime int64) {
	if m.tracker.Advance(currentTime) {
		m.triggerMap.Trigger(m.tracker.GetPosition())
	}
}
func (m *Manager) SetTrackerTime(currentTime int64) {
	m.tracker.SetCurrentTime(currentTime)
}

func (m *Manager) SetBPM(bpm float64) {
	m.tracker.SetBPM(bpm)
}

func (m *Manager) SetTrigger(pos BeatPosition, trigger func()) {
	m.triggerMap.AddTrigger(pos, trigger)
}

func (m *Manager) Clear() {
	m.triggerMap.Clear()
}

type BeatTriggerMap struct {
	triggers map[BeatPosition][]func()
}

func NewBeatTriggerMap() *BeatTriggerMap {
	return &BeatTriggerMap{
		triggers: make(map[BeatPosition][]func()),
	}
}

func (sm *BeatTriggerMap) AddTrigger(pos BeatPosition, trigger func()) {
	if _, exists := sm.triggers[pos]; !exists {
		sm.triggers[pos] = make([]func(), 0)
	}
	sm.triggers[pos] = append(sm.triggers[pos], trigger)
}

func (sm *BeatTriggerMap) Trigger(pos BeatPosition) {
	if triggers, exists := sm.triggers[pos]; exists {
		for _, trigger := range triggers {
			trigger()
		}
	}
}

func (sm *BeatTriggerMap) Clear() {
	sm.triggers = make(map[BeatPosition][]func())
}
