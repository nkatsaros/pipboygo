import ReactDOM from 'react-dom'
import React, { Component, PropTypes } from 'react'
import THREE from 'three'

const width = 768
const height = 360

import arrowImage from './arrow.gif'

export default class LocalMap extends Component {
  static propTypes = {
    buffer: PropTypes.object,
    color: PropTypes.array.isRequired,
    player: PropTypes.object.isRequired
  }

  componentDidMount() {
    this._scene = new THREE.Scene()
    this._camera = new THREE.OrthographicCamera(width / -2, width / 2, height / 2, height / -2, 1, 1000)

    this._renderer = new THREE.WebGLRenderer()

    this._renderer.setSize(width, height)
    this._container.appendChild(this._renderer.domElement)

    // setup map geometry
    const geometry = new THREE.PlaneGeometry(width, height)
    const material = new THREE.MeshBasicMaterial({ color: 0xffffff })
    this._map = new THREE.Mesh(geometry, material)
    this._scene.add(this._map)

    // setup arrow geometry
    const arrowGeom = new THREE.PlaneGeometry(32, 32)
    const arrowTexture = THREE.ImageUtils.loadTexture(arrowImage)
    arrowTexture.needsUpdate = true
    const arrowMat = new THREE.MeshBasicMaterial({ color: 0xffffff, map: arrowTexture, transparent: true })
    this._arrow = new THREE.Mesh(arrowGeom, arrowMat)
    this._scene.add(this._arrow)


    this._camera.position.z = 5

    this._renderer.render(this._scene, this._camera)
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.buffer !== this.props.buffer) {
      // dispose old texture
      if (this._map.material.map) {
        this._map.material.map.dispose()
      }

      const texture = new THREE.DataTexture(new Uint8Array(nextProps.buffer), width, height, THREE.LuminanceFormat, THREE.UnsignedByteType)
      texture.flipY = true
      texture.needsUpdate = true

      this._map.material.map = texture
      this._map.material.needsUpdate = true
    }

    if (nextProps.player !== this.props.player) {
      this._arrow.rotation.z = THREE.Math.degToRad(nextProps.player.Rotation)*-1
    }

    if (nextProps.color !== this.props.color) {
      // use selected pipboy color
      const color = new THREE.Color(nextProps.color[0], nextProps.color[1], nextProps.color[2])
      this._map.material.color.setHex(color.getHex())
      this._arrow.material.color.setHex(color.getHex())
      this._map.material.needsUpdate = true
      this._arrow.material.needsUpdate = true
    }

    if (nextProps !== this.props) {
      this._renderer.render(this._scene, this._camera)
    }
  }

  componentWillUnmount() {
  }

  shouldComponentUpdate() {
    return false
  }

  render() {
    return (
      <div ref={c => this._container = c} />
    )
  }
}
