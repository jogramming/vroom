package vroom

import (
	"github.com/veandco/go-sdl2/sdl"
)

type System interface {
	AddComponent(component Component)
	RemoveComponent(component Component)
	GetListenComponents() []string // The components this system will contain
	//GetComponents() []Component
	//ForEachComponent(cb func(Component) bool) // Callback takes component as arguement and returns wether to keep it in the slice or not
}

type BaseSystem struct {
	Components []Component
}

func (bs *BaseSystem) AddComponent(component Component) {
	bs.Components = append(bs.Components, component)
}

func (bs *BaseSystem) RemoveComponent(component Component) {
	index := -1
	for k, v := range bs.Components {
		if v == component {
			index = k
			break
		}
	}
	if index != -1 {
		if index == 0 {
			bs.Components = bs.Components[1:]
		} else if index == len(bs.Components)-1 {
			bs.Components = bs.Components[:len(bs.Components)-1]
		} else {
			bs.Components = append(bs.Components[:index], bs.Components[index+1:]...)
		}
	}
}

func (bs *BaseSystem) GetListenComponents() []string {
	return []string{""}
}

func (bs *BaseSystem) GetComponents() []Component {
	return nil
}

func (bs *BaseSystem) ClearComponents() {
	bs.Components = nil
}

func (bs *BaseSystem) ForEachComponent(cb func(Component) bool) {
	newSlice := make([]Component, 0)
	for _, v := range bs.Components {
		if v == nil {
			continue
		}

		if !v.GetParent().Enabled() {
			newSlice = append(newSlice, v)
			continue
		}

		keep := cb(v)
		if keep {
			newSlice = append(newSlice, v)
		}
	}
	bs.Components = newSlice
}

// Some core systems
type DrawSystem struct {
	components map[int][]*DrawComp
}

func (ds *DrawSystem) AddComponent(component Component) {
	cast, ok := component.(*DrawComp)
	if !ok {
		return
	}

	layer := cast.Layer

	if ds.components == nil {
		ds.components = make(map[int][]*DrawComp)
	}

	compSlice := ds.components[layer]
	compSlice = append(compSlice, cast)
	ds.components[layer] = compSlice
}

func (ds *DrawSystem) RemoveComponent(component Component) {
	cast, ok := component.(*DrawComp)
	if !ok {
		return
	}

	layer := cast.Layer

	compSlice := ds.components[layer]

	index := -1
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

	ds.components[layer] = compSlice
}

func (ds *DrawSystem) ClearComponents() {
	ds.components = nil
}

func (ds *DrawSystem) GetListenComponents() []string {
	return []string{"DrawComp"}
}

func (ds *DrawSystem) Draw(renderer *sdl.Renderer) {
	for i := -10; i < 10; i++ {
		compSlice := ds.components[i]

		for _, comp := range compSlice {
			if comp == nil {
				ds.RemoveComponent(comp)
				continue
			}
			if comp.Draw == nil {
				continue
			}

			if comp.GetParent().Enabled() {
				comp.Draw(renderer)
			}
		}
	}
}

type UpdateSystem struct {
	BaseSystem
}

func (us *UpdateSystem) GetListenComponents() []string {
	return []string{"UpdateComp"}
}

func (us *UpdateSystem) Update(dt float64) {
	us.ForEachComponent(func(comp Component) bool {
		cast, ok := comp.(*UpdateComp)
		if !ok {
			return false
		}

		if cast.Update == nil {
			return true
		}

		cast.Update(dt)
		return true
	})
}

type MouseClickSystem struct {
	BaseSystem
}

func (mc *MouseClickSystem) GetListenComponents() []string {
	return []string{"MouseClickComp"}
}

func (mc *MouseClickSystem) MouseButtonEvent(x, y, button int, up bool) {
	mc.ForEachComponent(func(comp Component) bool {
		cast, ok := comp.(*MouseClickComp)
		if !ok {
			return false
		}

		// Check if the callbacks are nil or not
		if up {
			if cast.MouseUp == nil {
				return true
			}
		} else {
			if cast.MouseDown == nil {
				return true
			}
		}

		mboxComp := cast.GetComponent("MouseBox")
		transformComp := cast.GetComponent("Transform")
		if mboxComp != nil && transformComp != nil {
			transform, ok := transformComp.(*Transform)
			mbox, ok2 := mboxComp.(*MouseBox)
			if ok && ok2 {
				if x > int(transform.Position.X) && x < int(transform.Position.X)+mbox.W &&
					y > int(transform.Position.Y) && y < int(transform.Position.Y)+mbox.H {
					if up {
						cast.MouseUp(x, y, button)
					} else {
						cast.MouseDown(x, y, button)
					}
				}
			}
		} else {
			if up {
				cast.MouseUp(x, y, button)
			} else {
				cast.MouseDown(x, y, button)
			}
		}
		return true
	})
}

type MouseHoverSystem struct {
	BaseSystem
}

func (mh *MouseHoverSystem) GetListenComponents() []string {
	return []string{"MouseHoverComp"}
}

// 	MouseeEnter func()
// 	MouseLeave func()
// 	MouseMove
func (mh *MouseHoverSystem) MouseMove(x, y int) {
	mh.ForEachComponent(func(comp Component) bool {
		cast, ok := comp.(*MouseHoverComp)
		if !ok {
			return false
		}
		mboxComp := cast.GetComponent("MouseBox")
		transformComp := cast.GetComponent("Transform")
		if mboxComp != nil && transformComp != nil {
			transform, ok := transformComp.(*Transform)
			mbox, ok2 := mboxComp.(*MouseBox)
			if ok && ok2 {
				if x > int(transform.Position.X) && x < int(transform.Position.X)+mbox.W &&
					y > int(transform.Position.Y) && y < int(transform.Position.Y)+mbox.H {
					if !mbox.Active {
						if cast.MouseEnter != nil {
							cast.MouseEnter()
						}
						mbox.Active = true
					}
					if cast.MouseMove != nil {
						cast.MouseMove(x, y)
					}
				} else {
					if mbox.Active {
						if cast.MouseLeave != nil {
							cast.MouseLeave()
						}
						mbox.Active = false
					}
				}
			}
		} else {
			cast.MouseMove(x, y)
		}
		return true
	})
}

type KeyboardSystem struct {
	BaseSystem
}

func (kb *KeyboardSystem) GetListenComponents() []string {
	return []string{"KeyboardComp"}
}

func (kb *KeyboardSystem) KeyboardEvent(key sdl.Keycode, up bool) {
	kb.ForEachComponent(func(comp Component) bool {
		casted, ok := comp.(*KeyboardComp)
		if !ok {
			return false
		}

		if up {
			if casted.KeyUp != nil {
				casted.KeyUp(key)
			}
		} else {
			if casted.KeyDown != nil {
				casted.KeyDown(key)
			}
		}
		return true
	})
}
