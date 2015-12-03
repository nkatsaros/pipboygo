import ReactDOM from 'react-dom'
import React, { Component } from 'react'
import { connect } from 'react-redux'

import LocalMap from '../components/LocalMap'
import { getValue, resolveAddress } from '../reducers/memory'

import '../index.css'

function mapStateToProps(state) {
  // TODO: Remove this and do something a lot better
  if (!state.info.loaded) {
    return {
      loading: true
    }
  }

  return {
    loading: false,
    localmap: state.localmap,
    effectColor: getValue(state, ["Status", "EffectColor"]),
    localPlayer: getValue(state, ["Map", "Local", "Player"])
  }
}

class App extends Component {
  render() {
    const { loading } = this.props;

    if (loading) {
      return <div>loading</div>
    }

    const { localmap, effectColor, localPlayer } = this.props;

    return (
      <div>
        <LocalMap buffer={localmap} color={effectColor} player={localPlayer} />
      </div>
    )
  }
}

export default connect(mapStateToProps)(App)
