<script>
  import GraphEvent from "$lib/components/GraphEvent.svelte";
  import Gantt from "$lib/components/Gantt.svelte";

  let networkData = $state({ events: [], activities: [] });
  let loading = $state(false);

  function generateKey() {
    return Math.random().toString(36).substr(2, 9);
  }

  // Domyślne dane zgodne z Twoją dokumentacją API
  let inputData = $state([
    {
      _key: generateKey(),
      id: "A",
      from_node: "1",
      to_node: "2",
      duration: "3",
    },
    {
      _key: generateKey(),
      id: "B",
      from_node: "1",
      to_node: "3",
      duration: "5",
    },
    {
      _key: generateKey(),
      id: "C",
      from_node: "2",
      to_node: "4",
      duration: "4",
    },
    {
      _key: generateKey(),
      id: "D",
      from_node: "3",
      to_node: "4",
      duration: "6",
    },
  ]);

  function addRow() {
    inputData.push({
      _key: generateKey(),
      id: "",
      from_node: "",
      to_node: "",
      duration: "1",
    });
  }

  /** @param {number} index */
  function removeRow(index) {
    inputData.splice(index, 1);
  }

  async function fetchActivities() {
    loading = true;

    // Przekształcamy dane do formatu oczekiwanego przez /event
    const payload = inputData
      .filter(
        (row) =>
          row.id.trim() !== "" && row.from_node !== "" && row.to_node !== "",
      )
      .map((row) => ({
        id: row.id.trim(),
        // Wymuszamy typ Integer (liczba całkowita) dla węzłów i czasu
        from_node: parseInt(row.from_node, 10),
        to_node: parseInt(row.to_node, 10),
        duration: parseInt(row.duration, 10) || 0,
      }));

    try {
      // Wywołanie poprawnego endpointu z dokumentacji
      const response = await fetch("http://localhost:8080/event", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!response.ok) throw new Error("Błąd sieci");
      networkData = await response.json();
    } catch (e) {
      console.error("Błąd pobierania danych", e);
      alert("Dane są nieprawidłowe.");
    } finally {
      loading = false;
    }
  }

  function exportToCSV() {
    const escapeCSV = (value) => {
      value = String(value || "");
      if (value.includes(",") || value.includes('"')) {
        return '"' + value.replace(/"/g, '""') + '"';
      }
      return value;
    };

    const headers = ["ID", "From Node", "To Node", "Duration"];
    const rows = inputData.map((row) => [
      escapeCSV(row.id),
      escapeCSV(row.from_node),
      escapeCSV(row.to_node),
      escapeCSV(row.duration),
    ]);
    const csvContent = [headers, ...rows]
      .map((row) => row.join(","))
      .join("\n");
    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = "event_tasks.csv";
    link.click();
  }

  function importFromCSV(event) {
    const parseCSVLine = (line) => {
      const result = [];
      let current = "";
      let inQuotes = false;
      for (let i = 0; i < line.length; i++) {
        const char = line[i];
        if (char === '"') {
          if (inQuotes && line[i + 1] === '"') {
            current += '"';
            i++; // skip next quote
          } else {
            inQuotes = !inQuotes;
          }
        } else if (char === "," && !inQuotes) {
          result.push(current);
          current = "";
        } else {
          current += char;
        }
      }
      result.push(current);
      return result;
    };

    const file = event.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      const text = e.target.result;
      const lines = text.trim().split("\n");
      if (lines.length < 2) {
        alert(
          "Plik CSV musi zawierać co najmniej nagłówek i jeden wiersz danych.",
        );
        return;
      }

      const headers = parseCSVLine(lines[0]).map((h) => h.trim());
      const expectedHeaders = ["ID", "From Node", "To Node", "Duration"];
      if (
        headers.length !== 4 ||
        !headers.every((h, i) => h === expectedHeaders[i])
      ) {
        alert(
          "Nieprawidłowy format pliku CSV. Oczekiwane nagłówki: ID,From Node,To Node,Duration",
        );
        return;
      }

      const parsedData = [];
      for (let i = 1; i < lines.length; i++) {
        const cols = parseCSVLine(lines[i]);
        if (cols.length !== 4) {
          alert(`Wiersz ${i + 1}: Nieprawidłowa liczba kolumn. Oczekiwane 4.`);
          return;
        }
        const id = cols[0].trim();
        const fromNode = cols[1].trim();
        const toNode = cols[2].trim();
        const durationStr = cols[3].trim();

        if (!id) {
          alert(`Wiersz ${i + 1}: ID zadania nie może być puste.`);
          return;
        }
        if (!fromNode || isNaN(parseInt(fromNode, 10))) {
          alert(`Wiersz ${i + 1}: Węzeł początkowy musi być liczbą całkowitą.`);
          return;
        }
        if (!toNode || isNaN(parseInt(toNode, 10))) {
          alert(`Wiersz ${i + 1}: Węzeł końcowy musi być liczbą całkowitą.`);
          return;
        }
        const duration = parseInt(durationStr, 10);
        if (isNaN(duration) || duration <= 0) {
          alert(
            `Wiersz ${i + 1}: Czas trwania musi być dodatnią liczbą całkowitą.`,
          );
          return;
        }

        parsedData.push({
          _key: generateKey(),
          id,
          from_node: fromNode,
          to_node: toNode,
          duration: durationStr,
        });
      }

      inputData = parsedData;
      networkData = { events: [], activities: [] }; // Wyczyść poprzednie wyniki obliczeń
      alert("Dane zostały pomyślnie zaimportowane.");
    };
    reader.readAsText(file);
  }

  /**
   * @param {any} data
   * @returns {any[]}
   */
  function buildGanttTasks(data) {
    const etByNode = new Map(
      data.events.map(
        /** @param {any} ev */ (ev) => [String(ev.node), Number(ev.et) || 0],
      ),
    );

    const predecessorsByNode = data.activities.reduce(
      /**
       * @param {Map<string, string[]>} map
       * @param {any} act
       */
      (map, act) => {
        const node = String(act.to_node);
        const previous = map.get(node) ?? [];
        map.set(node, [...previous, act.id]);
        return map;
      },
      new Map(),
    );

    return data.activities.map(
      /** @param {any} act */ (act) => ({
        id: act.id,
        duration: Number(act.duration) || 0,
        start: etByNode.get(String(act.from_node)) ?? 0,
        predecessors: predecessorsByNode.get(String(act.from_node)) ?? [],
        isCritical: act.is_critical ?? act.isCritical ?? false,
      }),
    );
  }

  const ganttTasks = $derived(() => buildGanttTasks(networkData));
</script>

<div class="layout">
  <a class="btn-add" href="/">Przełącz na czynnościowy</a>
  <div class="card form-section">
    <h2>Wariant zdarzeniowy</h2>

    <div class="table-header">
      <span>Zadanie</span>
      <span>Węzeł Początkowy</span>
      <span>Węzeł Końcowy</span>
      <span>Czas trwania</span>
      <span>Akcje</span>
    </div>

    {#each inputData as row, index (row._key)}
      <div class="table-row">
        <input
          type="text"
          bind:value={row.id}
          placeholder="np. A"
          class="input-field"
        />
        <input
          type="number"
          bind:value={row.from_node}
          placeholder="np. 1"
          class="input-field"
        />
        <input
          type="number"
          bind:value={row.to_node}
          placeholder="np. 2"
          class="input-field"
        />
        <input
          type="number"
          bind:value={row.duration}
          min="1"
          class="input-field"
        />
        <button
          class="btn-remove"
          onclick={() => removeRow(index)}
          aria-label="Usuń wiersz">✖</button
        >
      </div>
    {/each}

    <div class="form-actions">
      <button class="btn-add" onclick={addRow}>+ Dodaj czynność</button>
      <label class="btn-import">
        Importuj z CSV
        <input
          type="file"
          accept=".csv"
          onchange={importFromCSV}
          style="display: none;"
        />
      </label>
      <button class="btn-export" onclick={exportToCSV}>Eksportuj do CSV</button>
      <button class="btn-submit" onclick={fetchActivities} disabled={loading}>
        {loading ? "Obliczanie..." : "Oblicz ścieżkę krytyczną"}
      </button>
    </div>
  </div>

  {#if networkData.events.length > 0}
    <div class="card">
      <GraphEvent {networkData} />
    </div>

    {#if ganttTasks().length > 0}
      <div class="card">
        <Gantt activities={ganttTasks()} />
      </div>
    {/if}
  {/if}
</div>

<style>
  .layout {
    display: flex;
    flex-direction: column;
    gap: 24px;
    max-width: 900px;
    margin: 0 auto;
    padding: 20px;
    font-family: system-ui, sans-serif;
  }
  .card {
    background: #fff;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
  }
  h2 {
    margin-top: 0;
    font-size: 1.25rem;
    color: #333;
    border-bottom: 2px solid #f1f5f9;
    padding-bottom: 12px;
    margin-bottom: 20px;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }
  .btn-png {
    background: #10b981;
    color: white;
    padding: 8px 12px;
    font-size: 14px;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }
  .btn-png:hover {
    background: #059669;
  }
  .table-header,
  .table-row {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr 1fr 50px;
    gap: 12px;
    align-items: center;
  }
  .table-header {
    font-weight: 600;
    font-size: 0.85rem;
    color: #64748b;
    margin-bottom: 8px;
    padding: 0 4px;
  }
  .table-row {
    margin-bottom: 10px;
  }
  .input-field {
    width: 100%;
    padding: 10px;
    border: 1px solid #cbd5e1;
    border-radius: 6px;
    font-size: 0.95rem;
    box-sizing: border-box;
  }
  .input-field:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
  button {
    font-family: inherit;
    font-size: inherit;
    cursor: pointer;
    font-weight: 600;
    border: none;
    border-radius: 6px;
    transition: all 0.2s;
  }
  .btn-remove {
    background: #fee2e2;
    color: #ef4444;
    height: 38px;
    width: 38px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .btn-remove:hover {
    background: #fca5a5;
  }
  .form-actions {
    display: flex;
    justify-content: space-between;
    margin-top: 20px;
    padding-top: 16px;
    border-top: 1px solid #f1f5f9;
  }
  .btn-add {
    background: #f1f5f9;
    color: #475569;
    padding: 10px 16px;
  }
  .btn-add:hover {
    background: #e2e8f0;
  }
  .btn-import {
    background: #f59e0b;
    color: white;
    padding: 10px 16px;
    cursor: pointer;
    font-weight: 600;
    border: none;
    border-radius: 6px;
    transition: all 0.2s;
  }
  .btn-import:hover {
    background: #d97706;
  }
  .btn-export {
    background: #10b981;
    color: white;
    padding: 10px 16px;
  }
  .btn-export:hover {
    background: #059669;
  }
  .btn-submit {
    background: #3b82f6;
    color: white;
    padding: 10px 24px;
    font-size: 1rem;
  }
  .btn-submit:hover:not(:disabled) {
    background: #2563eb;
  }
  .btn-submit:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
</style>
