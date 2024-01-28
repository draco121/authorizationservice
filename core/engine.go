package core

import (
	"github.com/draco121/common/constants"
	"slices"
)

func authorizationEngine(allowedActions []constants.Action, requiredActions []constants.Action) bool {
	if slices.Contains(allowedActions, constants.All) {
		return true
	} else {
		for i := range allowedActions {
			if !slices.Contains(requiredActions, allowedActions[i]) {
				return false
			}
		}
		return true
	}
}
