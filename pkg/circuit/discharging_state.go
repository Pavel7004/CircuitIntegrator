package circuit

import (
	"math"
)

type dischargingState struct {
	circ *Circuit
}

var _ circuitState = (*dischargingState)(nil)

func newDischargingState(circ *Circuit) *dischargingState {
	return &dischargingState{
		circ: circ,
	}
}

func (s *dischargingState) GetDerivative() *Derivative {
	capVolts := make([]float64, 0, s.circ.stagesCount)
	for _, capVol := range s.circ.voltagesCap {
		capVolts = append(capVolts, -capVol/s.circ.load.tau)
	}
	return &Derivative{
		capVolts: capVolts,
	}
}

func (s *dischargingState) CheckDerivative(step float64, d *Derivative) bool {
	return s.circ.voltagesCap[0]+step*d.capVolts[0] > s.circ.holdingVoltage
}

func (s *dischargingState) CalculateOptimalStep(oldStep float64, d *Derivative) float64 {
	err := math.Abs(s.circ.holdingVoltage - s.circ.voltagesCap[0] + oldStep*d.capVolts[0])
	return 1 / err
}

func (s *dischargingState) Clone(newCirc *Circuit) circuitState {
	return &dischargingState{
		circ: newCirc,
	}
}

func (s *dischargingState) GetLoadVoltage() float64 {
	var capVoltage float64
	for _, vol := range s.circ.voltagesCap {
		capVoltage += vol
	}
	return capVoltage
}

func (s *dischargingState) ChangeState() {
	s.circ.state = newChargingState(s.circ)
}
