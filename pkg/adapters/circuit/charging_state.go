package circuit

type chargingState struct {
	circ *Circuit
}

var _ circuitState = (*chargingState)(nil)

func newChargingState(circ *Circuit) *chargingState {
	return &chargingState{
		circ: circ,
	}
}

func (s *chargingState) GetDerivative() *Derivative {
	capVolts := make([]float64, 0, s.circ.stagesCount)
	for i, volCap := range s.circ.voltagesCap {
		capVolts = append(capVolts, (s.circ.supplyVoltage-volCap)/s.circ.tau[i])
	}
	return &Derivative{
		capVolts: capVolts,
	}
}

func (s *chargingState) CheckDerivative(step float64, d *Derivative) bool {
	return s.circ.voltagesCap[0]+step*d.capVolts[0]-s.circ.gapTriggerVoltage < FloatPointAccuracy
}

func (s *chargingState) CalculateOptimalStep(oldStep float64, d *Derivative) float64 {
	l := 0.0
	r := oldStep

	for r-l > FloatPointAccuracy {
		m := (l + r) / 2

		if s.circ.voltagesCap[0]+d.capVolts[0]*m-s.circ.gapTriggerVoltage <= FloatPointAccuracy {
			l = m
		} else {
			r = m
		}
	}

	return l
}

func (s *chargingState) Clone(newCirc *Circuit) circuitState {
	return &chargingState{
		circ: newCirc,
	}
}

func (s *chargingState) GetLoadVoltage() float64 {
	return 0.0
}

func (s *chargingState) ChangeState() {
	s.circ.state = newDischargingState(s.circ)
}
