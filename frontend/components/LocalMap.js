import ReactDOM from 'react-dom'
import React, { Component, PropTypes } from 'react'

export default class LocalMap extends Component {
  static propTypes = {
    url: PropTypes.string.isRequired,
    r: PropTypes.number,
    g: PropTypes.number,
    b: PropTypes.number
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
    const { r, g, b } = this.props
    this._canvas.width = this._image.width
    this._canvas.height = this._image.height

    URL.revokeObjectURL(this._image.src)
    this._ctx.globalCompositeOperation = "source-over"
    this._ctx.drawImage(this._image, 0, 0)
    this._ctx.globalCompositeOperation = "multiply"
    this._ctx.fillStyle = `rgb(${Math.round((r||0)*255)}, ${Math.round((g||1)*255)}, ${Math.round((b||0)*255)})`
    this._ctx.fillRect(0, 0, this._canvas.width, this._canvas.height)
  }

  shouldComponentUpdate() {
    return false
  }

  render() {
    return (
      <canvas ref={c => this._canvas = c} />
    )
  }
}
