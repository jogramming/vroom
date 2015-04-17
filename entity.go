package vroom

type Entity interface {
	Init()  // Called when the entity is supposed to be initialized
	Start() // Called when the entity is added to the currently active scene

	InitCalled() bool // Returns wether init has been called or not
	SetInitCalled()   // Marks the init as called

	AddComponent(component Component)
	RemoveComponent(component Component)

	GetComponents() map[string][]Component
	GetComponentsByName(name string) []Component
	GetComponent(name string) Component // Get the first component by this name

	Enabled() bool
	SetEnabled(enable bool)

	GetEngine() *Engine
	SetEngine(e *Engine)

	GetParent() Entity
	SetParent(e Entity)

	AddChild(e Entity, addToScene bool)
	RemoveChild(e Entity, removeFromScene bool)
	GetChildren(recursive bool) []Entity

	Destroy()
}

type BaseEntity struct {
	Components map[string][]Component
	Entities   []Entity
	IsInit     bool
	Disabled   bool
	Engine     *Engine
	Parent     Entity
}

// Use the init functions to add the components
func (be *BaseEntity) Init() {}
func (be *BaseEntity) InitCalled() bool {
	return be.IsInit
}
func (be *BaseEntity) SetInitCalled() {
	be.IsInit = true
}

func (be *BaseEntity) Enabled() bool {
	return !be.Disabled
}

func (be *BaseEntity) SetEnabled(enable bool) {
	be.Disabled = !enable
}

func (be *BaseEntity) GetEngine() *Engine {
	return be.Engine
}

func (be *BaseEntity) SetEngine(e *Engine) {
	be.Engine = e
}

func (be *BaseEntity) Start() {}

func (be *BaseEntity) AddComponent(component Component) {
	if be.Components == nil {
		be.Components = make(map[string][]Component, 0)
	}
	compSlice, _ := be.Components[component.Name()]
	compSlice = append(compSlice, component)
	be.Components[component.Name()] = compSlice

	component.SetParent(be)
}
func (be *BaseEntity) RemoveComponent(component Component) {
	index := -1

	compSlice := be.Components[component.Name()]

	for k, v := range compSlice {
		if v == component {
			index = k
			break
		}
	}
	if index != -1 {
		if index == 0 {
			compSlice = compSlice[1:]
		} else if index == len(compSlice)-1 {
			compSlice = compSlice[:len(compSlice)-1]
		} else {
			compSlice = append(compSlice[:index], compSlice[index+1:]...)
		}
	}
	be.Components[component.Name()] = compSlice
}
func (be *BaseEntity) GetComponents() map[string][]Component {
	return be.Components
}
func (be *BaseEntity) GetComponentsByName(name string) []Component {
	return be.Components[name]
}
func (be *BaseEntity) GetComponent(name string) Component { // Get the first component by this name
	slice, _ := be.Components[name]
	if len(slice) < 1 {
		return nil
	}
	return slice[0]
}

func (be *BaseEntity) GetParent() Entity {
	return be.Parent
}

func (be *BaseEntity) SetParent(e Entity) {
	be.Parent = e
}

func (be *BaseEntity) AddChild(e Entity, addToScene bool) {
	be.Entities = append(be.Entities, e)
	if addToScene {
		be.GetEngine().AddEntity(e)
	}

	e.SetParent(be)
}

func (be *BaseEntity) RemoveChild(e Entity, removeFromScene bool) {
	index := -1
	for k, v := range be.Entities {
		if v == e {
			index = k
			break
		}
	}
	if index != -1 {
		if index == 0 {
			be.Entities = be.Entities[1:]
		} else if index == len(be.Entities)-1 {
			be.Entities = be.Entities[:len(be.Entities)-1]
		} else {
			be.Entities = append(be.Entities[:index], be.Entities[index+1:]...)
		}
	}

	if removeFromScene {
		be.GetEngine().RemoveEntity(e)
	}

	e.SetParent(nil)
}

func (be *BaseEntity) GetChildren(recursive bool) []Entity {
	entities := make([]Entity, 0)
	entities = append(entities, be.Entities...)

	if recursive {
		for _, ent := range entities {
			children := ent.GetChildren(true)
			if len(children) > 0 {
				entities = append(entities, children...)
			}
		}
	}
	return entities
}

// Destroy this entity and brutally murder all it's children
func (be *BaseEntity) Destroy() {
	for _, v := range be.Entities {
		v.Destroy()
	}
	for _, compSlice := range be.Components {
		for _, comp := range compSlice {
			comp.Destroy()
		}
	}
}

// Not actually empty contains a transform
func NewEntity(x, y int) Entity {
	ent := &BaseEntity{}
	ent.AddComponent(NewTransform(float64(x), float64(y), 0))
	return ent
}
