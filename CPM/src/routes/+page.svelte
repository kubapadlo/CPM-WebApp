<script>
  import Graph from "$lib/components/Graph.svelte";
  import Gantt from "$lib/components/Gantt.svelte";

  let activities = $state([]);
  let loading = $state(false);

  let graphComponent;
  let ganttComponent;

  function handleExport(event) {
    const { dataUrl, filename } = event.detail;
    const link = document.createElement("a");
    link.href = dataUrl;
    link.download = filename;
    link.click();
  }

  function generateKey() {
    return Math.random().toString(36).substr(2, 9);
  }

  let inputData = $state([
    { _key: generateKey(), id: "A", duration: 3, predecessors: "" },
    { _key: generateKey(), id: "B", duration: 5, predecessors: "A" },
    { _key: generateKey(), id: "C", duration: 2, predecessors: "A" },
    { _key: generateKey(), id: "D", duration: 4, predecessors: "B" },
    { _key: generateKey(), id: "E", duration: 6, predecessors: "C" },
    { _key: generateKey(), id: "F", duration: 3, predecessors: "B" },
    { _key: generateKey(), id: "G", duration: 5, predecessors: "D, E" },
    { _key: generateKey(), id: "H", duration: 2, predecessors: "F" },
    { _key: generateKey(), id: "I", duration: 4, predecessors: "G, H" },
    { _key: generateKey(), id: "J", duration: 3, predecessors: "I" },
    { _key: generateKey(), id: "K", duration: 2, predecessors: "J" },
  ]);

  function addRow() {
    inputData.push({
      _key: generateKey(),
      id: "",
      duration: 1,
      predecessors: "",
    });
  }

  function removeRow(index) {
    inputData.splice(index, 1);
  }

  async function fetchActivities() {
    loading = true;

    const payload = inputData
      .filter((row) => row.id.trim() !== "")
      .map((row) => ({
        id: row.id.trim(),
        duration: Number(row.duration) || 0,
        predecessors: row.predecessors
          .split(",")
          .map((p) => p.trim())
          .filter((p) => p !== ""),
      }));

    try {
      const response = await fetch("http://localhost:8080/activity", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      activities = await response.json();
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

    const headers = ["ID", "Duration", "Predecessors"];
    const rows = inputData.map((row) => [
      escapeCSV(row.id),
      escapeCSV(row.duration),
      escapeCSV(row.predecessors),
    ]);
    const csvContent = [headers, ...rows]
      .map((row) => row.join(","))
      .join("\n");
    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = "tasks.csv";
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
      const expectedHeaders = ["ID", "Duration", "Predecessors"];
      if (
        headers.length !== 3 ||
        !headers.every((h, i) => h === expectedHeaders[i])
      ) {
        alert(
          "Nieprawidłowy format pliku CSV. Oczekiwane nagłówki: ID,Duration,Predecessors",
        );
        return;
      }

      const parsedData = [];
      for (let i = 1; i < lines.length; i++) {
        const cols = parseCSVLine(lines[i]);
        if (cols.length !== 3) {
          alert(`Wiersz ${i + 1}: Nieprawidłowa liczba kolumn. Oczekiwane 3.`);
          return;
        }
        const id = cols[0].trim();
        const durationStr = cols[1].trim();
        const predecessors = cols[2].trim();

        if (!id) {
          alert(`Wiersz ${i + 1}: ID zadania nie może być puste.`);
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
          duration,
          predecessors,
        });
      }

      inputData = parsedData;
      activities = []; // Wyczyść poprzednie wyniki obliczeń
      alert("Dane zostały pomyślnie zaimportowane.");
    };
    reader.readAsText(file);
  }
</script>

<div class="layout">
  <a class="btn-add" href="/event">Przełącz na zdarzeniowy</a>
  <div class="card form-section">
    <h2>Wariant czynnościowy</h2>
    <div class="table-header">
      <span>Zadanie</span>
      <span>Czas trwania</span>
      <span>Poprzedniki (po przecinku)</span>
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
          bind:value={row.duration}
          min="1"
          class="input-field"
        />
        <input
          type="text"
          bind:value={row.predecessors}
          placeholder="np. A, B"
          class="input-field"
        />
        <button
          class="btn-remove"
          onclick={() => removeRow(index)}
          aria-label="Usuń wiersz"
        >
          ✖
        </button>
      </div>
    {/each}

    <div class="form-actions">
      <button class="btn-add" onclick={addRow}>+ Dodaj kolejne zadanie</button>
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

  {#if activities.length > 0}
    <div class="card">
      <div class="card-header">
        <h2>Graf</h2>
        <button
          class="btn-png"
          onclick={() => graphComponent?.()}
          title="Eksportuj wykres do PNG">📷 PNG</button
        >
      </div>
      <Graph bind:exportFunction={graphComponent} {activities} />
    </div>

    <div class="card">
      <div class="card-header">
        <h2>Wykres Gantta</h2>
        <button
          class="btn-png"
          onclick={() => ganttComponent?.()}
          title="Eksportuj wykres Gantta do PNG">📷 PNG</button
        >
      </div>
      <Gantt bind:exportFunction={ganttComponent} {activities} />
    </div>
  {/if}
</div>

<style>
  .layout {
    display: flex;
    flex-direction: column;
    gap: 24px;
    max-width: 1000px;
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

  .table-header,
  .table-row {
    display: grid;
    grid-template-columns: 1fr 1fr 2fr 50px;
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
    transition: border-color 0.2s;
  }

  .input-field:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  button {
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
  }

  .btn-png:hover {
    background: #059669;
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

  .btn-export {
    background: #10b981;
    color: white;
    padding: 10px 16px;
  }

  .btn-export:hover {
    background: #059669;
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

  pre {
    background: #f8fafc;
    padding: 16px;
    border-radius: 6px;
    overflow-x: auto;
    font-size: 0.85rem;
    border: 1px solid #e2e8f0;
  }
</style>
