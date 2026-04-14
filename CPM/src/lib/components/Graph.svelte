<script>
  import { onMount } from "svelte";
  import cytoscape from "cytoscape";
  // @ts-ignore
  import dagre from "cytoscape-dagre";
  import { createEventDispatcher } from "svelte";

  let { activities = [] } = $props();

  let container;
  let cy;
  let exportFunction = exportPNG;
  const dispatch = createEventDispatcher();

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
      dispatch("export", { dataUrl: pngData, filename: "graph.png" });
    }
  }
</script>

<div bind:this={container} class="graph-container"></div>

<style>
  .graph-container {
    width: 100%;
    height: 500px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background: #f9f9f9;
  }
</style>
