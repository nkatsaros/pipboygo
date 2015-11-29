import { combineReducers } from 'redux'
import database from './database'
import localmap from './localmap'

const rootReducer = combineReducers({
  database,
  localmap
})

export default rootReducer
