package circuit

type circuitState interface {
	GetDerivative() *Derivative
	Clone(newCirc *Circuit) circuitState
	GetLoadVoltage() float64
	ChangeState()
}
