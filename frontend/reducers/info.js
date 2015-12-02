import {
  MEMORY_UPDATE
} from '../actions'

export default function localmap(state = {
  loaded: false
}, action) {
  switch (action.type) {
    case MEMORY_UPDATE:
      return { ...state, loaded: true }
    default:
      return state
  }
}
