package circuit

type circuitState interface {
	GetDerivative() *Derivative
	CheckDerivative(step float64, d *Derivative) bool
	CalculateOptimalStep(oldStep float64, d *Derivative) float64
	Clone(newCirc *Circuit) circuitState
	GetLoadVoltage() float64
	ChangeState()
}
