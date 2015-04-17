package vroom

type Scene struct {
	Entities []Entity
}

func (s *Scene) AddEntity(entity Entity) {
	s.Entities = append(s.Entities, entity)
	if !entity.InitCalled() {
		entity.Init()
		entity.SetInitCalled()
	}
}
