import { combineReducers } from 'redux'
import memory from './memory'
import localmap from './localmap'
import info from './info'

const rootReducer = combineReducers({
  info,
  memory,
  localmap
})

export default rootReducer
