package circuit

type circuitState interface {
	GetDerivative() []float64
	Clone(newCirc *Circuit) circuitState
	GetLoadVoltage() float64
	ChangeState()
}
