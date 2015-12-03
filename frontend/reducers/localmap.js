import {
  LOCAL_MAP_UPDATE
} from '../actions'

export default function localmap(state = null, action) {
  switch (action.type) {
    case LOCAL_MAP_UPDATE:
      return action.buffer
    default:
      return state
  }
}
