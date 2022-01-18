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

func (s *chargingState) GetDerivative() []float64 {
	var (
		derivative          = make([]float64, 0, s.circ.stagesCount)
		prevCapVolt float64 = s.circ.supplyVoltage
	)
	for _, volCap := range s.circ.voltagesCap {
		chargeRatio := prevCapVolt / s.circ.supplyVoltage
		derivative = append(derivative, chargeRatio*(s.circ.supplyVoltage-volCap)/s.circ.tau)
		prevCapVolt = volCap
	}
	return derivative
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
