package circuit

type Derivative struct {
	capVolts []float64
}

func (d *Derivative) WeighCopy(weight float64) *Derivative {
	newCapVolts := make([]float64, 0, len(d.capVolts))
	for _, capVolt := range d.capVolts {
		newCapVolts = append(newCapVolts, weight*capVolt)
	}
	return &Derivative{
		capVolts: newCapVolts,
	}
}

func (d *Derivative) Add(weight float64, add *Derivative) *Derivative {
	for i := range d.capVolts {
		d.capVolts[i] += weight * add.capVolts[i]
	}
	return d
}
