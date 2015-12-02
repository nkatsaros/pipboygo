import {
  MEMORY_UPDATE
} from '../actions'

export default function memory(state = {}, action) {
  switch (action.type) {
    case MEMORY_UPDATE:
      let nextState = {...state, ...action.memory.added}
      for (let removed in action.memory.removed) {
        console.log("removed", removed)
        delete nextState[removed]
      }
      return nextState
    default:
      return state
  }
}

export function getValue(state, parts) {
  const mem = state.memory;
  if (!mem.hasOwnProperty('0')) {
    return null;
  }
  let currentValue = mem['0']
  for (let part of parts) {
    if (!currentValue.hasOwnProperty(part)) {
      throw "not found!"
    }

    currentValue = mem[currentValue[part]]
  }
  return currentValue
}
