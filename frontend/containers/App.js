import ReactDOM from 'react-dom'
import React, { Component } from 'react'
import { connect } from 'react-redux'

import LocalMap from '../components/LocalMap'
import { getValue } from '../reducers/database'

import '../index.css'

function mapStateToProps(state) {
  return {
    localmap: state.localmap,
    database: state.database,
    effectRed: getValue(state, ["Status", "EffectColor", 0]),
    effectGreen: getValue(state, ["Status", "EffectColor", 1]),
    effectBlue: getValue(state, ["Status", "EffectColor", 2]),
  }
}

class App extends Component {
  render() {
    const { database, localmap, effectRed, effectGreen, effectBlue } = this.props;

    return (
      <div>
        <LocalMap url={localmap} r={effectRed} g={effectGreen} b={effectBlue} />
        <p>{JSON.stringify(database)}</p>
      </div>
    )
  }
}

export default connect(mapStateToProps)(App)
