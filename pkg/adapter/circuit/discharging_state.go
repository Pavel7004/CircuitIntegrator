package circuit

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
	return s.circ.voltagesCap[0]+step*d.capVolts[0]-s.circ.holdingVoltage > FloatPointAccuracy
}

func (s *dischargingState) CalculateOptimalStep(oldStep float64, d *Derivative) float64 {
	l := 0.0
	r := oldStep

	for r-l > FloatPointAccuracy {
		m := (l + r) / 2

		if s.circ.voltagesCap[0]+d.capVolts[0]*m-s.circ.holdingVoltage <= FloatPointAccuracy {
			l = m
		} else {
			r = m
		}
	}

	if l < FloatPointAccuracy {
		return r
	}
	return l
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
