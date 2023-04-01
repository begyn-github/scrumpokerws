package statemachine

import (
	"reflect"
)

type UserState struct {
	Data        UserData
	ActualState *State
}

type UserData struct {
	UserName   string
	UserEmail  string
	SessionId  string
	TaskName   string
	StorePoint string
}

func (us *UserState) UpdateDataValue(value string) {
	r := reflect.ValueOf(us)
	f := reflect.Indirect(r).FieldByName("Data").FieldByName(us.ActualState.Field)

	f.SetString(value)
}
