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

func (s *dischargingState) GetDerivative() []float64 {
	var (
		derivative = make([]float64, 0, s.circ.stagesCount)
	)
	for _, capVol := range s.circ.voltagesCap {
		derivative = append(derivative, -capVol/s.circ.load.tau)
	}
	return derivative
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
