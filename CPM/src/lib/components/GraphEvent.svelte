<script>
  import cytoscape from "cytoscape";
  // @ts-ignore
  import dagre from "cytoscape-dagre";

  let { networkData = { events: [], activities: [] } } = $props();

  let container;
  let cy;

  $effect(() => {
    if (
      networkData &&
      (networkData.events.length > 0 || networkData.activities.length > 0) &&
      container
    ) {
      drawGraph();
    }
  });

  function drawGraph() {
    cytoscape.use(dagre);

    const elements = [];

    networkData.events.forEach((ev) => {
      elements.push({
        data: {
          id: ev.node.toString(),
          label: `${ev.node}\n──────\nET: ${ev.et.toString().padEnd(2)}\nLT: ${ev.lt.toString().padEnd(2)}`,
          isCritical: ev.is_critical,
        },
      });
    });

    networkData.activities.forEach((act) => {
      elements.push({
        data: {
          id: act.id,
          source: act.from_node.toString(),
          target: act.to_node.toString(),
          label: `${act.id}, Czas: ${act.duration}`,
          isCritical: act.is_critical,
        },
      });
    });

    cy = cytoscape({
      container: container,
      elements: elements,
      style: [
        {
          selector: "node",
          style: {
            label: "data(label)",
            "text-wrap": "wrap",
            "text-valign": "center",
            "font-size": "12px",
            "background-color": "#fff",
            "border-width": 2,
            "border-color": "#333",
            shape: "rectangle",
            padding: "10px",
            width: "label",
            height: "label",
          },
        },
        {
          selector: "edge",
          style: {
            label: "data(label)",
            "curve-style": "bezier",
            "target-arrow-shape": "triangle",
            width: 2,
            "line-color": "#ccc",
            "font-size": "12px",
            "text-rotation": "autorotate",
            "text-margin-y": -12,
            "text-background-color": "#f9f9f9",
            "text-background-opacity": 0.8,
          },
        },
        {
          selector: "node[?isCritical]",
          style: { "background-color": "#ff4d4f", color: "#fff" },
        },
      ],
      layout: { name: "dagre", rankDir: "LR", rankSep: 100 },
    });
  }

  function exportPNG() {
    if (cy) {
      const pngData = cy.png({ full: true, scale: 2 });
      const link = document.createElement("a");
      link.href = pngData;
      link.download = "graph_event.png";
      link.click();
    }
  }
</script>

<div class="graph-wrapper">
  <div class="graph-header">
    <h3>Graf zdarzeniowy</h3>
    <button class="btn-png" onclick={exportPNG} title="Eksportuj wykres do PNG"
      >📷 PNG</button
    >
  </div>
  <div bind:this={container} class="graph-container"></div>
</div>

<style>
  .graph-wrapper {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .graph-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .graph-header h3 {
    margin: 0;
  }

  .btn-png {
    background: #3b82f6;
    color: white;
    border: none;
    padding: 8px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
  }

  .btn-png:hover {
    background: #2563eb;
  }

  .graph-container {
    width: 100%;
    height: 500px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background: #f9f9f9;
  }
</style>
