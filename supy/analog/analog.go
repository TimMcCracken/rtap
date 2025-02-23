package analog

import (
	"fmt"
)

type AnalogLimits struct {
	HighReasonability float64
	HighCritical      float64
	HighOperating     float64
	HighWarning       float64
	LowWarning        float64
	LowOperating      float64
	LowCritical       float64
	LowReasonaility   float64
}

type AnalogClass struct {
	Limits        int64 //index into the []AnalogLimits
	ZeroDeadband  float64
	AlarmDeadband float64
}

type AnalogT struct {
}

type AnalogL struct {
	HighReasonability bool
	HighCritical      bool
	HighOperating     bool
	HighWarning       bool
	LowWarning        bool
	LowOperating      bool
	LowCritical       bool
	LowReasonaility   bool
}

type AnalogQ struct {
	dummy string // replace this later
}

type AnalogConfig struct {
	name  [16]byte
	class int64 // index into the []AnalogClass
}

type AnalogPoint struct {
	Name  [16]byte
	Value float64
	T     AnalogT
	L     AnalogL
	Q     AnalogQ
}

func (tag *AnalogPoint) Update(dq *AnalogQ, value float64) {

	fmt.Println("Hello from analog_tag", value)

}
