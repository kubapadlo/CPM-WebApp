<script>
  import { onMount } from "svelte";
  import cytoscape from "cytoscape";
  // @ts-ignore
  import dagre from "cytoscape-dagre";

  let { activities = [] } = $props();

  let container;
  let cy;

  $effect(() => {
    if (activities && container) {
      drawGraph();
    }
  });

  function drawGraph() {
    cytoscape.use(dagre);

    const elements = [];
    activities.forEach((act) => {
      elements.push({
        data: {
          id: act.id,
          label: `${act.id}\n──────────\nCzas: ${act.duration}\n Rezerwa: ${act.slack}\n──────────\nES: ${act.es} | EF: ${act.ef}\nLS: ${act.ls} | LF: ${act.lf}`,
          isCritical: act.is_critical,
        },
      });
      act.predecessors.forEach((pred) => {
        elements.push({
          data: {
            id: `${pred}-${act.id}`,
            source: pred,
            target: act.id,
            isCritical: act.is_critical,
          },
        });
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
            "curve-style": "bezier",
            "target-arrow-shape": "triangle",
            width: 2,
            "line-color": "#ccc",
          },
        },
        {
          selector: "[?isCritical]",
          style: { "background-color": "#ff4d4f", color: "#fff" },
        },
      ],
      layout: { name: "dagre", rankDir: "LR" },
    });
  }

  function exportPNG() {
    if (cy) {
      const pngData = cy.png({ full: true, scale: 2 });
      const link = document.createElement("a");
      link.href = pngData;
      link.download = "graph.png";
      link.click();
    }
  }
</script>

<div class="graph-wrapper">
  <div class="graph-header">
    <h3>Graf</h3>
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
