package configpropertyutil

type ConfigPropertiesUtil interface {
	GetProp(name string) string
	GetIntProp(name string) int
	GetInt32Prop(name string) int32
	GetInt64Prop(name string) int64
	GetBoolProp(name string) bool
	GetFloat64Prop(name string) float64
	GetStringArray(name string, splitBy string) []string
}
