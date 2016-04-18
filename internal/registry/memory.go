package registry

type MemoryRegistry struct {
	T string
}

func NewInMemRegistry() *MemoryRegistry {
	return &MemoryRegistry{T: "inmem"}
}

func (m *MemoryRegistry) Register(reg Registration) error {
	return nil
}

func (m *MemoryRegistry) DeRegister(id string) error {
	return nil
}

func (m *MemoryRegistry) Type() string {
	return m.T
}
