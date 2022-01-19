package circuit

import (
	"math"
)

type Circuit struct {
	supplyVoltage     float64
	capacity          float64
	resistance        float64
	stagesCount       uint
	gapTriggerVoltage float64
	load              LoadComponents
	state             circuitState
	tau               float64
	stateChange       bool
	voltagesCap       []float64
}

func NewCircuit(chargeComp ChargeComponents, load LoadComponents) *Circuit {
	load.chargeCapacity = chargeComp.Capacity / float64(chargeComp.StagesCount)
	load.tau = load.chargeCapacity * (load.Resistance + chargeComp.Resistance)
	circ := &Circuit{
		supplyVoltage:     chargeComp.SupplyVoltage,
		capacity:          chargeComp.Capacity,
		resistance:        chargeComp.Resistance,
		stagesCount:       chargeComp.StagesCount,
		gapTriggerVoltage: chargeComp.GapTriggerVoltage,
		load:              load,
		state:             nil,
		tau:               2 * chargeComp.Resistance * chargeComp.Capacity,
		stateChange:       true,
		voltagesCap:       make([]float64, chargeComp.StagesCount),
	}
	circ.state = newChargingState(circ)
	return circ
}

func (st *Circuit) GetDerivative() []float64 {
	return st.state.GetDerivative()
}

func (st *Circuit) ToggleState() {
	st.state.ChangeState()
}

func (st *Circuit) ApplyDerivative(h float64, derivative []float64) {
	for i := range st.voltagesCap {
		st.voltagesCap[i] += h * derivative[i]
		if st.voltagesCap[i] < 0 {
			st.voltagesCap[i] = 0
		}
	}
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
	return -st.tau * math.Log(1-st.gapTriggerVoltage/st.supplyVoltage)
}

func (st *Circuit) GetLoadVoltageFunc() func(x float64) float64 {
	stateChangeTime := st.GetSystemPeriod()
	return func(x float64) float64 {
		if x < stateChangeTime {
			return 0.0
		}
		return st.gapTriggerVoltage * float64(st.stagesCount) * math.Exp(-(x-stateChangeTime)/st.load.tau)
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
		load:              load,
		state:             nil,
		tau:               st.tau,
		stateChange:       false,
		voltagesCap:       make([]float64, 0, st.stagesCount),
	}
	for _, capVol := range st.voltagesCap {
		newCirc.voltagesCap = append(newCirc.voltagesCap, capVol)
	}
	newState := st.state.Clone(newCirc)
	newCirc.state = newState
	return newCirc
}
