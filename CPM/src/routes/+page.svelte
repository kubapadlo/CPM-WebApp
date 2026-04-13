<script>
  import Graph from '$lib/components/Graph.svelte';

  let activities = $state([]);
  let loading = $state(false);

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
    { _key: generateKey(), id: "K", duration: 2, predecessors: "J" }
  ]);

  function addRow() {
    inputData.push({ _key: generateKey(), id: "", duration: 1, predecessors: "" });
  }

  function removeRow(index) {
    inputData.splice(index, 1);
  }

  async function fetchActivities() {
    loading = true;

    const payload = inputData
      .filter(row => row.id.trim() !== '')
      .map(row => ({
        id: row.id.trim(),
        duration: Number(row.duration) || 0,
        predecessors: row.predecessors
          .split(',')
          .map(p => p.trim())
          .filter(p => p !== '') 
      }));

    try {
      const response = await fetch('http://localhost:8080/activity', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      activities = await response.json();
    } catch (e) {
      console.error("Błąd pobierania danych", e);
      alert("Wystąpił błąd podczas łączenia z serwerem.");
    } finally {
      loading = false;
    }
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
        <input type="text" bind:value={row.id} placeholder="np. A" class="input-field" />
        <input type="number" bind:value={row.duration} min="1" class="input-field" />
        <input type="text" bind:value={row.predecessors} placeholder="np. A, B" class="input-field" />
        <button class="btn-remove" onclick={() => removeRow(index)} aria-label="Usuń wiersz">
          ✖
        </button>
      </div>
    {/each}

    <div class="form-actions">
      <button class="btn-add" onclick={addRow}>+ Dodaj kolejne zadanie</button>
      <button class="btn-submit" onclick={fetchActivities} disabled={loading}>
        {loading ? 'Obliczanie...' : 'Oblicz ścieżkę krytyczną'}
      </button>
    </div>
  </div>

  {#if activities.length > 0}
    <div class="card">
      <h2>Graf</h2>
      <Graph {activities} />
      
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

  .table-header, .table-row {
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