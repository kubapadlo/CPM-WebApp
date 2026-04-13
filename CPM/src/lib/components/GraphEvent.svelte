<script>
  import { onMount } from "svelte";
  import cytoscape from "cytoscape";
  // @ts-ignore
  import dagre from "cytoscape-dagre";
  import { createEventDispatcher } from "svelte";

  let { networkData = { events: [], activities: [] } } = $props();

  let container;
  let cy;
  const dispatch = createEventDispatcher();

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

    // Dodajemy zdarzenia (WĘZŁY / KÓŁKA)
    networkData.events.forEach((ev) => {
      elements.push({
        data: {
          // Cytoscape WYMAGA stringów dla id, więc zamieniamy liczbę na tekst
          id: ev.node.toString(),
          // Używamy ET i LT zgodnie z backendem
          label: `${ev.node}\n──────\nET: ${ev.et.toString().padEnd(2)}\nLT: ${ev.lt.toString().padEnd(2)}`,
          isCritical: ev.is_critical,
        },
      });
    });

    // Dodajemy czynności (KRAWĘDZIE / STRZAŁKI)
    networkData.activities.forEach((act) => {
      elements.push({
        data: {
          id: act.id,
          // Łączymy strzałkę po numerach zdarzeń (również zamienionych na tekst)
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
    if (cy && exportFunction) {
      const pngData = cy.png({ full: true, scale: 2 });
      exportFunction(pngData, "graph_event.png");
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
