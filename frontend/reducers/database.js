import {
  DATABASE_UPDATE
} from '../actions'

export default function database(state = {}, action) {
  switch (action.type) {
    case DATABASE_UPDATE:
      return {...state, ...action.database}
    default:
      return state
  }
}

export function getValue(state, parts) {
  const db = state.database;
  if (!db.hasOwnProperty('0')) {
    return null;
  }
  let currentValue = db['0']
  for (let part of parts) {
    if (!currentValue.hasOwnProperty(part)) {
      throw "not found!"
    }

    currentValue = db[currentValue[part]]
  }
  return currentValue
}
