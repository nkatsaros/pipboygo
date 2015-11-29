import {
  LOCAL_MAP_UPDATE
} from '../actions'

export default function localmap(state = "", action) {
  switch (action.type) {
    case LOCAL_MAP_UPDATE:
      return action.url
    default:
      return state
  }
}
