package domain

var (
	ErrInvalidData = NewError(404, "invalid_data", "Error in received data")
)

type CircuitQuery struct {
	SupplyVoltage     float64 `json:"supply_voltage"`
	Capacity          float64 `json:"capacity"`
	Resistance        float64 `json:"resistance"`
	StagesCount       uint    `json:"stages_count"`
	GapTriggerVoltage float64 `json:"gap_trigger_voltage"`
	HoldingVoltage    float64 `json:"holding_voltage"`
	LoadResistance    float64 `json:"load_resistance"`
	Step              float64 `json:"step"`
	IntNum            int     `json:"int_num"`
}

func (q *CircuitQuery) Check() error {
	if q.IntNum < 0 && q.IntNum >= len(Integrators) {
		return ErrInvalidData
	}
	if q.Resistance < 0 {
		return ErrInvalidData
	}
	if q.Capacity < 0 {
		return ErrInvalidData
	}

	return nil
}
