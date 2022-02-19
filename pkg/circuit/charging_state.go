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

func (s *chargingState) Clone(newCirc *Circuit) circuitState {
	return &chargingState{
		circ: newCirc,
	}
}

func (s *chargingState) GetLoadVoltage() float64 {
	return 0.0
}

func (s *chargingState) ChangeState() {
	if s.circ.gapTriggerVoltage-s.circ.voltagesCap[0] < 0.0001 {
		s.circ.state = newDischargingState(s.circ)
	}
}
