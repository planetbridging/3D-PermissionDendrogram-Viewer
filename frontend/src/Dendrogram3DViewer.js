import React, { Component } from "react";
import { Canvas } from "@react-three/fiber";
import { OrbitControls, Html } from "@react-three/drei";
import * as THREE from "three";

// Color mapping based on overlap levels
const colors = {
  0: "gray", // Default
  1: "purple", // Overlapping only
  2: "yellow", // Hierarchical only
  3: "red", // Both overlapping and hierarchical
};

// Component to render each node as a sphere with a label
class SphereNode extends Component {
  render() {
    const { position, overlapLevel, name } = this.props;
    const color = colors[overlapLevel] || colors[0];

    return (
      <mesh position={position}>
        <sphereGeometry args={[0.5, 8, 8]} />
        <meshStandardMaterial color={color} />
        <Html
          position={[0, 1, 0]}
          center
          style={{ color: "white", fontSize: "0.5em" }}
        >
          {name}
        </Html>
      </mesh>
    );
  }
}

// Component to render a line connecting two points
class Line extends Component {
  constructor(props) {
    super(props);
    this.ref = React.createRef();
  }

  componentDidMount() {
    const { start, end } = this.props;
    const geometry = new THREE.BufferGeometry().setFromPoints([start, end]);
    this.ref.current.geometry = geometry;
  }

  componentDidUpdate() {
    const { start, end } = this.props;
    const geometry = new THREE.BufferGeometry().setFromPoints([start, end]);
    this.ref.current.geometry = geometry;
  }

  render() {
    return (
      <line ref={this.ref}>
        <bufferGeometry />
        <lineBasicMaterial color="gray" />
      </line>
    );
  }
}

// Main recursive component to render each node and its children
class DendrogramNode extends Component {
  render() {
    const { node, position, parentPosition, level } = this.props;
    const childNodes = node.children || [];

    return (
      <>
        {/* Render the current node as a sphere */}
        <SphereNode
          position={position}
          overlapLevel={node.overlapLevel}
          name={node.name}
        />

        {/* Draw a line to connect to the parent node, if it exists */}
        {parentPosition && <Line start={parentPosition} end={position} />}

        {/* Recursively render child nodes */}
        {childNodes.map((child, index) => {
          const angle = (2 * Math.PI * index) / childNodes.length;
          const radius = 5 + level * 2;
          const childPosition = new THREE.Vector3(
            position.x + Math.cos(angle) * radius,
            position.y - 4,
            position.z + Math.sin(angle) * radius
          );

          return (
            <DendrogramNode
              key={child.name + index}
              node={child}
              position={childPosition}
              parentPosition={position}
              level={level + 1}
            />
          );
        })}
      </>
    );
  }
}

// Main viewer component to set up the 3D scene and initiate the root node
class Dendrogram3DViewer extends Component {
  render() {
    return (
      <Canvas
        style={{ height: "100vh", background: "black" }}
        camera={{ position: [0, 0, 50], fov: 75 }}
      >
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} />
        {/* Render the root node */}
        <DendrogramNode
          node={this.props.dendrogram}
          position={new THREE.Vector3(0, 10, 0)}
          level={0}
        />
        {/* Controls for zoom, pan, and rotation */}
        <OrbitControls enableZoom={true} enableRotate={true} />
      </Canvas>
    );
  }
}

export default Dendrogram3DViewer;
