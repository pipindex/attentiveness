package Models

import (
	"fmt"
)

func NewAttentiveness(totalPoles, activePoles, idlePoles, hiddenPoles int) (a Attentiveness) {
	//equal to the points not found
	notPresentPoles := totalPoles - activePoles - idlePoles - hiddenPoles

	a = Attentiveness{
		ActiveRatio:     float64(activePoles) / float64(totalPoles),
		IdleRatio:       float64(idlePoles) / float64(totalPoles),
		HiddenRatio:     float64(hiddenPoles) / float64(totalPoles),
		NotPresentRatio: float64(notPresentPoles) / float64(totalPoles),
	}
	return a
}

type Attentiveness struct {
	ActiveRatio     float64
	IdleRatio       float64
	HiddenRatio     float64
	NotPresentRatio float64
}

func (a Attentiveness) String() string {
	return fmt.Sprintf("Active: %f%% Idle: %f%% Hidden: %f%% NotPresent: %f%% Total: %f%%",
		a.ActiveRatio*100, a.IdleRatio*100, a.HiddenRatio*100, a.NotPresentRatio*100,
		(a.ActiveRatio+a.IdleRatio+a.HiddenRatio+a.NotPresentRatio)*100)
}

func (a Attentiveness) ToMap() map[string]float64 {
	return map[string]float64{
		"active":      a.ActiveRatio,
		"idle":        a.IdleRatio,
		"hidden":      a.HiddenRatio,
		"not_present": a.NotPresentRatio,
	}
}
