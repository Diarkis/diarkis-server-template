package metrics

/*
ValueUpdate is a function that updates the associated metrics value by using the returned value as its associated metrics value.
*/
type ValueUpdate func() (value int)

/*
ValueLabelUpdate is a function that updates the associated metrics value by using the returned value as its associated metrics label and value.
*/
type ValueLabelUpdate func() (value []*LabelValue)

/*
LabelValue map[map[string]string]int
*/
type LabelValue struct {
	Value  interface{}       `json:"v"`
	Labels map[string]string `json:"l"`
}

/*
Item [INTERNAL USE ONLY] represents a metrics value
*/
type Item struct {
	Addr   string
	Value  interface{}
	Type   string
	Labels map[string]string
}

/*
Setup must be invoked at the start of the process in order to use metrics.

	[IMPORTANT] Must invoke this function BEFORE calling diarkis.Start()
	[IMPORTANT] MUST invoke this function AFTER calling mesh.Setup(...)
*/
func Setup(path string) {
}

/*
PassMeshOnUpdate [INTERNAL USE ONLY]
*/
func PassMeshOnUpdate(cb func(func())) {
}

/*
PassMeshGetNodeAddressesByRole [INTERNAL USE ONLY]
*/
func PassMeshGetNodeAddressesByRole(cb func(string) []string) {
}

/*
PassMeshGetNodeTypeByAddress [INTERNAL USE ONLY]
*/
func PassMeshGetNodeTypeByAddress(cb func(string) string) {
}

/*
PassMeshGetNodeValue [INTERNAL USE ONLY]
*/
func PassMeshGetNodeValue(cb func(string, string) interface{}) {
}

/*
PassMeshSetNodeValue [INTERNAL USE ONLY]
*/
func PassMeshSetNodeValue(cb func(string, interface{})) {
}

/*
DefineMetricsValue defines the value update operation of the metrics name given.

	[IMPORTANT] MUST invoke this function BEFORE calling diarkis.Start()
	[IMPORTANT] MUST invoke this function AFTER calling mesh.Setup(...)

Error Cases

	┌────────────────────────────────────────────────┬──────────────────────────────────────────────────────────────────────────────┐
	│ Error                                          │ Reason                                                                       │
	├────────────────────────────────────────────────┼──────────────────────────────────────────────────────────────────────────────┤
	│ Setup must invoked at the start of the process │ Invoking this function without calling Setup() is not allowed.               │
	│ Must be invoked BEFORE the start of Diarkis    │ Invoking this function after the start of Diarkis is not allowed.            │
	│ Metrics name must be configured                │ All metrics names must be configured in the configuration file.              │
	│ Metrics value update function already assigned │ Assigning multiple value update functions for a metrics name is not allowed. │
	└────────────────────────────────────────────────┴──────────────────────────────────────────────────────────────────────────────┘

The name must be configured in a configuration file.
*/
func DefineMetricsValue(name string, cb ValueUpdate) error {
	return nil
}

/*
DefineMetricsValueWithLabel defines the value update operation of the metrics name and label given.
Currently it only works for Prometheus metrics.
See the documentation for DefineMetricsValue for more information
*/
func DefineMetricsValueWithLabel(name string, cb ValueLabelUpdate) error {
	return nil
}

/*
GetMetricsValues [INTERNAL USE ONLY] returns all metrics values with the configured names and update functions
*/
func GetMetricsValues(nodeType string) map[string][]*Item {
	return nil
}
