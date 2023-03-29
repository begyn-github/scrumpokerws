package statemachine

type UserState struct {
	User        string
	ActualState *State
}
