package schema

import "math/big"

var Float0 = new(Float).SetUint64(0)

// Float is a wrapper for big decimal numbers
type Float struct {
	raw big.Float
}

func (z *Float) SetBigInt(b *big.Int) *Float {
	z.raw.SetString(b.String())
	return z
}

func (z *Float) SetString(s string) bool {
	_, ok := z.raw.SetString(s)
	return ok
}

func (z *Float) SetUint64(i uint64) *Float {
	z.raw.SetUint64(i)
	return z
}

func (z *Float) String() string {
	return z.raw.String()
}

func (f *Float) Shift(num int) *Float {
	str := "1"
	for i := 0; i < num; i++ {
		str += "0"
	}
	decF := new(Float)
	decF.SetString(str)
	return f.Div(decF)
}

func (f *Float) DivUint(i uint64) *Float {
	ii := new(Float).SetUint64(i)
	return f.Div(ii)
}

func (f *Float) Sub(j *Float) *Float {
	res := new(big.Float).Sub(&f.raw, &j.raw)
	return &Float{raw: *res}
}

func (f *Float) Add(j *Float) *Float {
	res := new(big.Float).Add(&f.raw, &j.raw)
	return &Float{raw: *res}
}

func (f *Float) Div(j *Float) *Float {
	res := new(big.Float).Quo(&f.raw, &j.raw)
	return &Float{raw: *res}
}
