package input

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	KEYBOARD_INPUT_TYPE = 1
	GAMEPAD_INPUT_TYPE = 2
	MOUSE_INPUT_TYPE = 3
)

type ActionHandler interface {
    IsActionJustPressed(string) bool
    IsActionJustReleased(string) bool
    ActionPressedDuration(string) int
}

type inputDef struct {
	input_type int

	keyboard_key ebiten.Key

	gamepad_id ebiten.GamepadID
	gamepad_key ebiten.GamepadButton

	mouse_key ebiten.MouseButton
}

type InputActionHandler struct {
	defined_actions map[string]inputDef
}

func NewInputActionHandler() {
	new := InputActionHandler{}
	new.defined_actions = make(map[string]inputDef)
}

func (i *InputActionHandler) RegisterKeyboardAction(action string, key ebiten.Key) {
	i.defined_actions[action] = inputDef{
		input_type: KEYBOARD_INPUT_TYPE,
		keyboard_key: key,
	}
}

func (i *InputActionHandler) RegisterGamepadAction(action string, id ebiten.GamepadID, key ebiten.GamepadButton) {
	i.defined_actions[action] = inputDef{
		input_type: GAMEPAD_INPUT_TYPE,
		gamepad_key: key,
		gamepad_id: id,
	}
}

func (i *InputActionHandler) RegisterMouseAction(action string, key ebiten.MouseButton) {
	i.defined_actions[action] = inputDef{
		input_type: MOUSE_INPUT_TYPE,
		mouse_key: key,
	}
}

func (i *InputActionHandler) UnregisterAction(action string) {
	delete(i.defined_actions, action)
}

func (i InputActionHandler) IsActionJustPressed(action string) bool {
	def := i.defined_actions[action]
	switch def.input_type {
	case KEYBOARD_INPUT_TYPE:
		return inpututil.IsKeyJustPressed(def.keyboard_key)
	case GAMEPAD_INPUT_TYPE:
		return inpututil.IsGamepadButtonJustPressed(def.gamepad_id, def.gamepad_key)
	case MOUSE_INPUT_TYPE:
		return inpututil.IsMouseButtonJustPressed(def.mouse_key)
	default:
		return false
	}

	return false
}

func (i InputActionHandler) IsActionJustReleased(action string) bool {
	def := i.defined_actions[action]
	switch def.input_type {
	case KEYBOARD_INPUT_TYPE:
		return inpututil.IsKeyJustReleased(def.keyboard_key)
	case GAMEPAD_INPUT_TYPE:
		return inpututil.IsGamepadButtonJustReleased(def.gamepad_id, def.gamepad_key)
	case MOUSE_INPUT_TYPE:
		return inpututil.IsMouseButtonJustReleased(def.mouse_key)
	default:
		return false
	}

	return false
}

func (i InputActionHandler) ActionPressedDuration(action string) int {
	def := i.defined_actions[action]
	switch def.input_type {
	case KEYBOARD_INPUT_TYPE:
		return inpututil.KeyPressDuration(def.keyboard_key)
	case GAMEPAD_INPUT_TYPE:
		return inpututil.GamepadButtonPressDuration(def.gamepad_id, def.gamepad_key)
	case MOUSE_INPUT_TYPE:
		return inpututil.MouseButtonPressDuration(def.mouse_key)
	default:
		return 0
	}

	return 0
}
