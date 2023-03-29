package statemachine

import "errors"

type State struct {
	Name string
	Id   int
	Near []Transition
}

type Transition struct {
	Word    string
	StateTo *State
}

var rootState State

func GetRoot() State {
	return rootState
}

func (s State) GetNear() []Transition {
	return s.Near
}

func (s State) GetAvailableWords() map[string]string {
	var availableWords = make(map[string]string)

	for _, tr := range s.GetNear() {
		availableWords[tr.Word] = tr.StateTo.Name
	}

	return availableWords
}

func (s State) GoTo(word string) (State, error) {
	var state *State = nil

	for _, tr := range s.GetNear() {
		if tr.Word == word {
			state = tr.StateTo
			break
		}
	}

	if state == nil {
		return State{}, errors.New("State not found")
	}

	return *state, nil
}

func (s State) GetMenu() []string {
	var output []string

	for key, value := range s.GetAvailableWords() {
		output = append(output, "["+key+"] - "+value)
	}

	return output
}

func init() {
	root,
		takeUserName,
		takeUserEmail,
		createSession,
		loginSession,
		createTask,
		choiceNumber,
		exitState1 :=
		State{Name: "Root", Id: 1},
		State{Name: "Set User Name", Id: 2},
		State{Name: "Set User Email", Id: 3},
		State{Name: "Create New Session", Id: 4},
		State{Name: "Login Session", Id: 5},
		State{Name: "Create New Task", Id: 6},
		State{Name: "Choice Store Point", Id: 7},
		State{Name: "Exit", Id: 8}

	root.Near = []Transition{{Word: "N", StateTo: &takeUserName}, {Word: "X", StateTo: &exitState1}}
	takeUserName.Near = []Transition{{Word: "M", StateTo: &takeUserEmail}, {Word: "X", StateTo: &exitState1}}
	takeUserEmail.Near = []Transition{{Word: "S", StateTo: &createSession}, {Word: "L", StateTo: &loginSession}, {Word: "X", StateTo: &exitState1}}
	createSession.Near = []Transition{{Word: "T", StateTo: &createTask}, {Word: "X", StateTo: &exitState1}}
	createTask.Near = []Transition{{Word: "P", StateTo: &choiceNumber}, {Word: "X", StateTo: &exitState1}}
	loginSession.Near = []Transition{{Word: "P", StateTo: &choiceNumber}, {Word: "X", StateTo: &exitState1}}
	choiceNumber.Near = []Transition{{Word: "P", StateTo: &choiceNumber}, {Word: "X", StateTo: &exitState1}}

	rootState = root
}
