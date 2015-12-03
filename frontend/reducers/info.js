import {
  MEMORY_UPDATE
} from '../actions'

export default function localmap(state = {
  loaded: false
}, action) {
  switch (action.type) {
    case MEMORY_UPDATE:
      // check is addr 0 was deleted
      let loading = action.memory.removed.filter(removed => removed === 0).length > 0
      return { ...state, loaded: !loading }
    default:
      return state
  }
}
