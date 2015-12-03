import {
  MEMORY_UPDATE
} from '../actions'

export default function memory(state = {}, action) {
  switch (action.type) {
    case MEMORY_UPDATE:
      let nextState = {...state, ...action.memory.added}
      action.memory.removed.forEach(removed => delete nextState[removed])
      return nextState
    default:
      return state
  }
}

export function getValue(state, parts) {
  const memory = state.memory
  if (!memory.hasOwnProperty('0')) {
    return null
  }

  let currentValue = memory['0']

  for (let i = 0; i < parts.length; i++) {
    const part = parts[i]

    if (!currentValue.hasOwnProperty(part)) {
      throw "not found!"
    }

    // found what we're looking for, resolveAddress it
    if (i == parts.length-1) {
      return resolveAddress(state, currentValue[part])
    }

    // keep iterating
    currentValue = memory[currentValue[part]]
  }

  return currentValue
}

export function resolveAddress(state, addr) {
  const memory = state.memory

  if (!memory.hasOwnProperty(addr)) {
    throw "not found!"
  }

  const val = memory[addr]

  if (Array.isArray(val)) {
    // arrays
    return val.map(addr => resolveAddress(state, addr))
  } else if (typeof val === 'object') {
    // dictionaries
    let res = {}
    for (let prop in val) {
      if (val.hasOwnProperty(prop)) {
        res[prop] = resolveAddress(state, val[prop])
      }
    }
    return res
  }
  return val
}
