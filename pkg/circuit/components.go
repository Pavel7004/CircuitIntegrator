package circuit

type ChargeComponents struct {
	SupplyVoltage     float64
	Capacity          float64
	Resistance        float64
	StagesCount       uint
	GapTriggerVoltage float64
}

type LoadComponents struct {
	Resistance     float64
	chargeCapacity float64
	tau            float64
}
