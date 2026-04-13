<script>
  import html2canvas from "html2canvas";
  import { createEventDispatcher } from "svelte";

  let { activities = [] } = $props();

  const dispatch = createEventDispatcher();

  /**
   * @param {any} item
   * @returns {string[]}
   */
  function normalizePredecessors(item) {
    if (Array.isArray(item.predecessors)) return item.predecessors;
    if (!item.predecessors) return [];
    return item.predecessors
      .toString()
      .split(",")
      .map(trimString)
      .filter(notEmptyString);
  }

  /**
   * @param {string} value
   * @returns {string}
   */
  function trimString(value) {
    return value.trim();
  }

  /**
   * @param {string} value
   * @returns {boolean}
   */
  function notEmptyString(value) {
    return value !== "";
  }

  /**
   * @param {any[]} items
   * @returns {any[]}
   */
  function buildSchedule(items) {
    /** @type {any[]} */
    const nodes = items.map((item) => ({
      id: item.id,
      duration: Number(item.duration) || 0,
      predecessors: normalizePredecessors(item),
      explicitStart: item.start ?? item.es ?? item.earliestStart ?? null,
      isCritical: item.is_critical ?? item.isCritical ?? false,
    }));

    const nodeMap = new Map(nodes.map((node) => [node.id, node]));
    const visited = new Set();
    const temp = new Set();
    /** @type {any[]} */
    const order = [];

    /**
     * @param {string} id
     */
    function visit(id) {
      if (visited.has(id)) return;
      if (temp.has(id)) return;
      temp.add(id);
      const node = nodeMap.get(id);
      if (!node) {
        temp.delete(id);
        return;
      }
      node.predecessors.forEach(
        /** @param {string} pred */ (pred) => visit(pred),
      );
      temp.delete(id);
      visited.add(id);
      order.push(node);
    }

    nodeMap.forEach((_, id) => visit(id));

    const schedule = new Map();
    order.forEach((node) => {
      /**
       * @param {number} max
       * @param {string} pred
       * @returns {number}
       */
      function reducer(max, pred) {
        const predEntry = schedule.get(pred);
        return Math.max(max, predEntry?.end ?? 0);
      }

      const predecessorEnd = node.predecessors.reduce(reducer, 0);
      const start =
        node.explicitStart !== null
          ? Math.max(Number(node.explicitStart), predecessorEnd)
          : predecessorEnd;
      schedule.set(node.id, {
        ...node,
        start,
        end: start + node.duration,
      });
    });

    return Array.from(schedule.values());
  }

  const tasks = $derived(() => buildSchedule(activities));
  const maxEnd = $derived(() =>
    tasks().length ? Math.max(...tasks().map((task) => task.end)) : 0,
  );
  const timeScale = 40;

  /** @param {number} max */
  function generateTimeLabels(max) {
    const labels = [];
    for (let i = 0; i <= max; i++) {
      labels.push(i);
    }
    return labels;
  }

  const timeLabels = $derived(() => generateTimeLabels(maxEnd()));
  const chartWidth = $derived(() => (maxEnd() + 1) * timeScale);

  let ganttContainer;
  let exportFunction = exportPNG;

  function exportPNG() {
    if (ganttContainer) {
      html2canvas(ganttContainer, {
        backgroundColor: "#ffffff",
        scale: 2,
        useCORS: true,
        allowTaint: false,
      })
        .then((canvas) => {
          const dataUrl = canvas.toDataURL("image/png");
          dispatch("export", { dataUrl, filename: "gantt.png" });
        })
        .catch((err) => {
          console.error("Błąd eksportu PNG:", err);
          alert("Wystąpił błąd podczas eksportu wykresu.");
        });
    }
  }
</script>

<div class="gantt" bind:this={ganttContainer}>
  <div class="gantt-grid" style="position: relative;">
    <!-- Oś czasu z etykietami co trzecią jednostkę -->
    <div class="time-axis">
      {#each timeLabels() as label}
        <div
          class="time-label"
          class:first={label === 0}
          style="left: {label * timeScale}px"
        >
          {label}
        </div>
      {/each}
    </div>

    <!-- Linie pionowe dla każdej etykiety czasu -->
    <svg
      class="time-lines"
      width={chartWidth()}
      height="100%"
      style="position: absolute; top: 0; left: 0; pointer-events: none;"
    >
      {#each timeLabels() as label}
        <line
          x1={label * timeScale}
          y1="30"
          x2={label * timeScale}
          y2="100%"
          stroke="#64748b"
          stroke-width="1"
          stroke-dasharray="2,2"
        />
      {/each}
    </svg>

    {#each tasks() as task}
      <div class="task-row">
        <div class="task-bar-wrapper" style="min-width: {chartWidth()}px">
          <div
            class="task-bar {task.isCritical ? 'critical' : ''}"
            style="left: {task.start * timeScale}px; width: {task.duration *
              timeScale -
              24}px"
            title="{task.id} — start: {task.start}, czas: {task.duration}, koniec: {task.end}"
          >
            {task.id}
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .gantt {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .gantt-grid {
    overflow-x: auto;
    position: relative;
  }

  .time-axis {
    position: relative;
    height: 30px;
    margin-bottom: 8px;
  }

  .time-label {
    position: absolute;
    top: 0;
    transform: translateX(-50%);
    font-size: 0.85rem;
    font-weight: 600;
    color: #475569;
    background: white;
    padding: 2px 6px;
    border-radius: 4px;
    border: 1px solid #e2e8f0;
    white-space: nowrap;
  }

  .time-label.first {
    transform: none;
    left: 0 !important;
  }

  .time-lines {
    z-index: 1;
  }

  .task-row {
    position: relative;
    min-height: 44px;
  }

  .task-bar-wrapper {
    position: relative;
    min-height: 44px;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    background: #f8fafc;
    overflow: hidden;
  }

  .task-bar {
    position: absolute;
    top: 6px;
    bottom: 6px;
    border-radius: 999px;
    background: linear-gradient(135deg, #3b82f6, #60a5fa);
    color: white;
    padding: 8px 12px;
    white-space: nowrap;
    font-size: 0.9rem;
    box-shadow: 0 2px 8px rgba(15, 23, 42, 0.18);
  }

  .task-bar.critical {
    background: linear-gradient(135deg, #ef4444, #f87171);
    box-shadow: 0 2px 12px rgba(220, 38, 38, 0.35);
  }
</style>
