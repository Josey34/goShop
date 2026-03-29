package valueobject

type Money int64

func NewMoney(cents int64) Money {
	return Money(cents)
}

func (m Money) Add(other Money) Money {
	return Money(m.Value() + other.Value())
}

func (m Money) Subtract(other Money) Money {
	return Money(m.Value() - other.Value())
}

func (m Money) Multiply(factor int) Money {
	return Money(m.Value() * int64(factor))
}

func (m Money) Percentage(pct int) Money {
	return Money(int64(float64(m.Value()) * (float64(pct) / 100)))
}

func (m Money) Value() int64 {
	return int64(m)
}

func (m Money) IsNegative() bool {
	return m.Value() < 0
}

func (m Money) IsZero() bool {
	return m.Value() == 0
}
