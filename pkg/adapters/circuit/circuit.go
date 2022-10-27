package circuit

import (
	"math"
)

const FloatPointAccuracy = 1e-10

type Circuit struct {
	supplyVoltage     float64
	capacity          float64
	resistance        float64
	stagesCount       uint
	gapTriggerVoltage float64
	holdingVoltage    float64
	load              LoadComponents
	tolerance         float64
	state             circuitState
	tau               []float64
	voltagesCap       []float64
}

func NewCircuit(chargeComp ChargeComponents, load LoadComponents) *Circuit {
	load.chargeCapacity = chargeComp.Capacity / float64(chargeComp.StagesCount)
	load.tau = load.chargeCapacity * load.Resistance
	circ := &Circuit{
		supplyVoltage:     chargeComp.SupplyVoltage,
		capacity:          chargeComp.Capacity,
		resistance:        chargeComp.Resistance,
		stagesCount:       chargeComp.StagesCount,
		gapTriggerVoltage: chargeComp.GapTriggerVoltage,
		holdingVoltage:    chargeComp.HoldingVoltage,
		load:              load,
		tolerance:         0.0001,
		state:             nil,
		tau:               make([]float64, chargeComp.StagesCount),
		voltagesCap:       make([]float64, chargeComp.StagesCount),
	}
	for i := range circ.tau {
		circ.tau[i] = 2 * float64(i+1) * circ.resistance * circ.capacity
	}
	circ.state = newChargingState(circ)
	return circ
}

func (st *Circuit) GetDerivative() *Derivative {
	return st.state.GetDerivative()
}

func (st *Circuit) ToggleState() {
	st.state.ChangeState()
}

func (st *Circuit) CheckDerivative(step float64, d *Derivative) bool {
	return st.state.CheckDerivative(step, d)
}

func (st *Circuit) CalculateOptimalStep(oldStep float64, d *Derivative) float64 {
	return st.state.CalculateOptimalStep(oldStep, d)
}

func (st *Circuit) ApplyDerivative(h float64, d *Derivative) *Circuit {
	for i := range st.voltagesCap {
		st.voltagesCap[i] += h * d.capVolts[i]
	}
	return st
}

func (st *Circuit) GetCapVoltage(pos uint) float64 {
	if pos == 0 || pos > st.stagesCount {
		panic("Incorrect capacitor position")
	}
	return st.voltagesCap[pos-1]
}

func (st *Circuit) GetLoadVoltage() float64 {
	return st.state.GetLoadVoltage()
}

func (st *Circuit) GetSystemPeriod() float64 {
	return -st.tau[0] * math.Log(1-st.gapTriggerVoltage/st.supplyVoltage)
}

func (st *Circuit) GetLoadVoltageFunc() func(x float64) float64 {
	var (
		stateChangeTime = st.GetSystemPeriod()
		maxVoltage      = st.supplyVoltage * float64(st.stagesCount)
	)
	for i := range st.tau {
		maxVoltage -= st.supplyVoltage * math.Exp(-stateChangeTime/st.tau[i])
	}
	return func(x float64) float64 {
		if x <= stateChangeTime {
			return 0.0
		}
		return maxVoltage * math.Exp(-(x-stateChangeTime)/st.load.tau)
	}
}

func (st *Circuit) GetSystemCurrent() float64 {
	var current float64
	for _, capVol := range st.voltagesCap {
		current += (st.supplyVoltage - capVol) / (2 * st.resistance)
	}
	return current
}

func (st *Circuit) Clone() *Circuit {
	load := LoadComponents{
		Resistance:     st.load.Resistance,
		chargeCapacity: st.load.chargeCapacity,
		tau:            st.load.tau,
	}
	newCirc := &Circuit{
		supplyVoltage:     st.supplyVoltage,
		capacity:          st.capacity,
		resistance:        st.resistance,
		stagesCount:       st.stagesCount,
		gapTriggerVoltage: st.gapTriggerVoltage,
		holdingVoltage:    st.holdingVoltage,
		load:              load,
		tolerance:         st.tolerance,
		state:             nil,
		tau:               st.tau,
		voltagesCap:       make([]float64, 0, st.stagesCount),
	}
	newCirc.voltagesCap = append(newCirc.voltagesCap, st.voltagesCap...)
	newState := st.state.Clone(newCirc)
	newCirc.state = newState
	return newCirc
}
