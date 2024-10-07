package task

// LabelSet manages labels assigned to a task definition.
type LabelSet struct {
	rawLabels map[string]interface{}
}

// Construct the LabelSet with required fields.
func NewLabelSet(labelOpts ...LabelOpt) *LabelSet {
	labelSet := &LabelSet{
		rawLabels: map[string]interface{}{},
	}
	for _, lo := range labelOpts {
		lo.Write(labelSet)
	}
	return labelSet
}

func (l *LabelSet) Get(key string) (any, bool) {
	value, exist := l.rawLabels[key]
	return value, exist
}

func (l *LabelSet) GetOrDefault(key string, defaultValue any) any {
	value, exist := l.Get(key)
	if exist {
		return value
	}
	return defaultValue
}

func (l *LabelSet) Set(key string, value any) {
	l.rawLabels[key] = value
}

// LabelOpt is a utility to modify label during task definition generation.
// See producer.go for the actual usage.
type LabelOpt interface {
	Write(label *LabelSet)
}

type labelAdder struct {
	key   string
	value any
}

// Write implements LabelOpt.
func (a *labelAdder) Write(label *LabelSet) {
	label.Set(a.key, a.value)
}

var _ LabelOpt = (*labelAdder)(nil)

func WithLabel(key string, value any) LabelOpt {
	return &labelAdder{
		key:   key,
		value: value,
	}
}

type cloneLabel struct {
	source *LabelSet
}

// Write implements LabelOpt.
func (c *cloneLabel) Write(label *LabelSet) {
	for key, value := range c.source.rawLabels {
		label.Set(key, value)
	}
}

var _ LabelOpt = (*cloneLabel)(nil)

func FromLabelSet(labelSet *LabelSet) LabelOpt {
	return &cloneLabel{
		source: labelSet,
	}
}

type LabelFilter interface {
	Filter(label *LabelSet) bool
}

type equalLabelFilter struct {
	labelKey      string
	defaultResult bool
	operand       any
}

// Filter implements LabelFilter.
func (e *equalLabelFilter) Filter(label *LabelSet) bool {
	l, exist := label.Get(e.labelKey)
	if exist {
		return l == e.operand
	}
	return e.defaultResult
}

var _ LabelFilter = (*equalLabelFilter)(nil)

func EqualLabelFilter(labelKey string, operand any, defaultResult bool) LabelFilter {
	return &equalLabelFilter{
		labelKey:      labelKey,
		operand:       operand,
		defaultResult: defaultResult,
	}
}
