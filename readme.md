#Vroom
Simple game engine using SDL for graphics/input and chipmunk for physics

This is a simple game engine written using go, it is still very much work in progress, and things will change.

The goal is to have it somewhat ready before ludum dare 32 (as that is when i plan to use it)

##Core

###Core systems

####Update

Adds components that implements the updateable interface. Calls update every frame(with deltatime, in ms, argument)

####Draw

Drawable interface

####mouseclick

Adds components that implments the mouseclicklistener interface, If the parent entity also has the mbox and transform components, it will only listen for clicks in that box

####mousehover

Adds components that implements the mousehoverlistener interface, Acts the same way as mouseclick listener if mbox + transform is on parent entity

####keyboard

Adds components that implement the keyboardlistener interface

###Core components

####Button

A button

####Sprite

Sprite, displays a image

####Label

Renders text