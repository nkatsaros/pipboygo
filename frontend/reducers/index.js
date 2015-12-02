import { combineReducers } from 'redux'
import memory from './memory'
import localmap from './localmap'

const rootReducer = combineReducers({
  memory,
  localmap
})

export default rootReducer
