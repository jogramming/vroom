package vroom

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type System interface {
	AddComponent(component Component)
	RemoveComponent(component Component)
	CleanUp()
	LastCleanUp() time.Time
	Clear()
}

type BaseSystem struct {
	Components  []Component
	lastCleanUp time.Time
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
	return bs.Components
}

func (bs *BaseSystem) ClearComponents() {
	bs.Components = nil
}

// cleanup remove nil components
func (bs *BaseSystem) CleanUp() {
	bs.lastCleanUp = time.Now()
	newSlice := make([]Component, 0)
	for _, v := range bs.Components {
		if v == nil {
			continue
		}
		newSlice = append(newSlice, v)
	}
}

func (bs *BaseSystem) Clear() {
	bs.Components = make([]Component, 0)
}

func (bs *BaseSystem) ForEachComponent(cb func(Component) bool) {
	for _, v := range bs.Components {
		if v == nil {
			continue
		}

		if !v.GetParent().Enabled() {
			continue
		}

		cb(v)
	}
}

func (bs *BaseSystem) LastCleanUp() time.Time {
	return bs.lastCleanUp
}

// Some core systems
type DrawSystem struct {
	components  map[int][]DrawAble
	lastCleanUp time.Time
}

func (ds *DrawSystem) Clear() {
	ds.components = make(map[int][]DrawAble)
}

func (ds *DrawSystem) LastCleanUp() time.Time {
	return ds.lastCleanUp
}

// cleanup remove nil components, maybe not needed since we do not do any expensive checking here
func (ds *DrawSystem) CleanUp() {
	ds.lastCleanUp = time.Now()
	for k, layer := range ds.components {
		newSlice := make([]DrawAble, 0)
		for _, v := range layer {
			if v == nil {
				continue
			}
			newSlice = append(newSlice, v)
		}
		ds.components[k] = newSlice
	}
}

func (ds *DrawSystem) AddComponent(component Component) {
	cast, ok := component.(DrawAble)
	if !ok {
		return
	}

	layer := cast.GetLayer()

	if ds.components == nil {
		ds.components = make(map[int][]DrawAble)
	}

	compSlice := ds.components[layer]
	compSlice = append(compSlice, cast)
	ds.components[layer] = compSlice
}

func (ds *DrawSystem) RemoveComponent(component Component) {
	cast, ok := component.(DrawAble)
	if !ok {
		return
	}

	layer := cast.GetLayer()

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

func (ds *DrawSystem) Draw(renderer *sdl.Renderer) {
	for i := -10; i < 10; i++ {
		compSlice := ds.components[i]
		for _, comp := range compSlice {
			if comp == nil {
				ds.RemoveComponent(comp)
				continue
			}

			if comp.Enabled() && (comp.GetParent() != nil && comp.GetParent().Enabled()) {
				comp.Draw(renderer)
			}
		}
	}
}

type UpdateSystem struct {
	BaseSystem
}

func (us *UpdateSystem) AddComponent(component Component) {
	_, ok := component.(UpdateAble)
	if ok {
		if us.Components == nil {
			us.Components = make([]Component, 0)
		}
		us.Components = append(us.Components, component)
	}
}

func (us *UpdateSystem) Update(dt float64) {
	us.ForEachComponent(func(comp Component) bool {
		cast, ok := comp.(UpdateAble)
		if !ok {
			return false
		}

		cast.Update(dt)
		return true
	})
}

type MouseClickSystem struct {
	BaseSystem
}

func (mc *MouseClickSystem) AddComponent(component Component) {
	_, ok := component.(MouseClickListener)
	if ok {
		if mc.Components == nil {
			mc.Components = make([]Component, 0)
		}
		mc.Components = append(mc.Components, component)
	}
}

func (mc *MouseClickSystem) MouseButtonEvent(x, y, button int, up bool) {
	mc.ForEachComponent(func(comp Component) bool {
		cast, ok := comp.(MouseClickListener)
		if !ok {
			return false
		}

		// Check if the callbacks are nil or not
		mboxComp := cast.GetComponent("MouseBox")
		transformComp := cast.GetComponent("Transform")
		if mboxComp != nil && transformComp != nil {
			transform, ok := transformComp.(*Transform)
			mbox, ok2 := mboxComp.(*MouseBox)
			position := transform.Position
			position.X -= float64(mbox.W / 2)
			position.Y -= float64(mbox.H / 2)
			if ok && ok2 {
				if x > int(position.X) && x < int(position.X)+mbox.W &&
					y > int(position.Y) && y < int(position.Y)+mbox.H {
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

func (mh *MouseHoverSystem) AddComponent(component Component) {
	_, ok := component.(MouseHoverListener)
	if ok {
		if mh.Components == nil {
			mh.Components = make([]Component, 0)
		}
		mh.Components = append(mh.Components, component)
	}
}

func (mh *MouseHoverSystem) MouseMove(x, y int) {
	mh.ForEachComponent(func(comp Component) bool {
		cast, ok := comp.(MouseHoverListener)
		if !ok {
			return false
		}
		mboxComp := cast.GetComponent("MouseBox")
		transformComp := cast.GetComponent("Transform")
		if mboxComp != nil && transformComp != nil {
			transform, ok := transformComp.(*Transform)
			mbox, ok2 := mboxComp.(*MouseBox)
			position := transform.Position
			position.X -= float64(mbox.W / 2)
			position.Y -= float64(mbox.H / 2)
			if ok && ok2 {
				if x > int(position.X) && x < int(position.X)+mbox.W &&
					y > int(position.Y) && y < int(position.Y)+mbox.H {
					if !mbox.Active {
						cast.MouseEnter()
						mbox.Active = true
					}
					cast.MouseMove(x, y)
				} else {
					if mbox.Active {
						cast.MouseLeave()
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
	Keys map[sdl.Keycode]bool
}

func (kb *KeyboardSystem) AddComponent(component Component) {
	_, ok := component.(KeyboardListener)
	if ok {
		if kb.Components == nil {
			kb.Components = make([]Component, 0)
		}
		kb.Components = append(kb.Components, component)
	}
}

func (kb *KeyboardSystem) KeyboardEvent(key sdl.Keycode, up bool) {
	if kb.Keys == nil {
		kb.Keys = make(map[sdl.Keycode]bool)
	}

	kb.Keys[key] = !up

	kb.ForEachComponent(func(comp Component) bool {
		casted, ok := comp.(KeyboardListener)
		if !ok {
			return false
		}

		if up {
			casted.KeyUp(key)
		} else {
			casted.KeyDown(key)
		}
		return true
	})
}
