import React from "react";
import Dendrogram3DViewer from "./Dendrogram3DViewer";

class DendrogramClient extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      dendrogram: null,
      status: "Disconnected",
    };
  }

  componentDidMount() {
    this.connectWebSocket();
  }

  connectWebSocket() {
    const ws = new WebSocket("ws://localhost:8432/ws/admin");

    ws.onopen = () => {
      this.setState({ status: "Connected" });
      console.log("Connected to WebSocket as admin");
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("Received dendrogram data from server:", data);
      this.setState({ dendrogram: data });
    };

    ws.onclose = () => {
      this.setState({ status: "Disconnected" });
      console.log("WebSocket connection closed");
    };

    ws.onerror = (error) => {
      console.log("WebSocket error:", error);
    };
  }

  render() {
    return (
      <div>
        <h1>Admin Dendrogram 3D Viewer</h1>
        <p>Status: {this.state.status}</p>
        {this.state.dendrogram ? (
          <Dendrogram3DViewer dendrogram={this.state.dendrogram} />
        ) : (
          <p>Loading dendrogram data...</p>
        )}
      </div>
    );
  }
}

export default DendrogramClient;
