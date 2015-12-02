import ReactDOM from 'react-dom'
import React, { Component, PropTypes } from 'react'

import styles from './LocalMap.css'

export default class LocalMap extends Component {
  static propTypes = {
    url: PropTypes.string.isRequired,
    color: PropTypes.array.isRequired,
    player: PropTypes.object.isRequired
  }

  componentDidMount() {
    this._image = new Image()
    this._image.addEventListener('load', ::this.handleLoad)

    this._ctx = this._canvas.getContext('2d')
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.url !== this.props.url) {
      this._image.src = nextProps.url
    }
  }

  componentWillUnmount() {
    this._image.removeEventListener('load', ::this.handleLoad)
  }

  handleLoad() {
    const { color } = this.props
    this._canvas.width = this._image.width
    this._canvas.height = this._image.height

    URL.revokeObjectURL(this._image.src)
    this._ctx.globalCompositeOperation = "source-over"
    this._ctx.drawImage(this._image, 0, 0)
    this._ctx.globalCompositeOperation = "multiply"
    this._ctx.fillStyle = `rgb(${Math.round((color[0]||0)*255)}, ${Math.round((color[1]||1)*255)}, ${Math.round((color[2]||0)*255)})`
    this._ctx.fillRect(0, 0, this._canvas.width, this._canvas.height)
  }

  // shouldComponentUpdate() {
  //   return false
  // }

  render() {
    const { player } = this.props

    return (
      <div className={styles.container}>
        <img className={styles.arrow} style={{transform: `rotate(${player.Rotation}deg)`}} src={require('./arrow.png')} />
        <canvas ref={c => this._canvas = c} />
      </div>
    )
  }
}
