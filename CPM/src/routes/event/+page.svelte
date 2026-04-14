<script>
  import GraphEvent from '$lib/components/GraphEvent.svelte';

  let networkData = $state({ events: [], activities:[] });
  let loading = $state(false);

  function generateKey() {
    return Math.random().toString(36).substr(2, 9);
  }

  // Domyślne dane zgodne z Twoją dokumentacją API
  let inputData = $state([
    { _key: generateKey(), id: "A", from_node: 1, to_node: 2, duration: 3 },
    { _key: generateKey(), id: "B", from_node: 1, to_node: 3, duration: 5 },
    { _key: generateKey(), id: "C", from_node: 2, to_node: 4, duration: 4 },
    { _key: generateKey(), id: "D", from_node: 3, to_node: 4, duration: 6 }
  ]);

  function addRow() {
    inputData.push({ _key: generateKey(), id: "", from_node: "", to_node: "", duration: 1 });
  }

  function removeRow(index) {
    inputData.splice(index, 1);
  }

  async function fetchActivities() {
    loading = true;

    // Przekształcamy dane do formatu oczekiwanego przez /event
    const payload = inputData
      .filter(row => row.id.trim() !== '' && row.from_node !== '' && row.to_node !== '')
      .map(row => ({
        id: row.id.trim(),
        // Wymuszamy typ Integer (liczba całkowita) dla węzłów i czasu
        from_node: parseInt(row.from_node, 10),
        to_node: parseInt(row.to_node, 10),
        duration: parseInt(row.duration, 10) || 0
      }));

    try {
      // Wywołanie poprawnego endpointu z dokumentacji
      const response = await fetch('http://localhost:8080/event', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      
    //   if (!response.ok) throw new Error("Błąd sieci");
      networkData = await response.json();
      
    } catch (e) {
      console.error("Błąd pobierania danych", e);
      alert("Dane są nieprawidłowe.");
      
    } finally {
      loading = false;
    }
  }
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
        <input type="text" bind:value={row.id} placeholder="np. A" class="input-field" />
        <input type="number" bind:value={row.from_node} placeholder="np. 1" class="input-field" />
        <input type="number" bind:value={row.to_node} placeholder="np. 2" class="input-field" />
        <input type="number" bind:value={row.duration} min="1" class="input-field" />
        <button class="btn-remove" onclick={() => removeRow(index)} aria-label="Usuń wiersz">✖</button>
      </div>
    {/each}

    <div class="form-actions">
      <button class="btn-add" onclick={addRow}>+ Dodaj czynność</button>
      <button class="btn-submit" onclick={fetchActivities} disabled={loading}>
        {loading ? 'Obliczanie...' : 'Oblicz ścieżkę krytyczną'}
      </button>
    </div>
  </div>

  {#if networkData.events.length > 0}
    <div class="card">
      <h2>Graf</h2>
      <GraphEvent {networkData} />
      
    </div>
  {/if}
</div>

<style>
  .layout { display: flex; flex-direction: column; gap: 24px; max-width: 900px; margin: 0 auto; padding: 20px; font-family: system-ui, sans-serif; }
  .card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 24px; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05); }
  h2 { margin-top: 0; font-size: 1.25rem; color: #333; border-bottom: 2px solid #f1f5f9; padding-bottom: 12px; margin-bottom: 20px; }
  .table-header, .table-row { display: grid; grid-template-columns: 1fr 1fr 1fr 1fr 50px; gap: 12px; align-items: center; }
  .table-header { font-weight: 600; font-size: 0.85rem; color: #64748b; margin-bottom: 8px; padding: 0 4px; }
  .table-row { margin-bottom: 10px; }
  .input-field { width: 100%; padding: 10px; border: 1px solid #cbd5e1; border-radius: 6px; font-size: 0.95rem; box-sizing: border-box; }
  .input-field:focus { outline: none; border-color: #3b82f6; box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1); }
  button { cursor: pointer; font-weight: 600; border: none; border-radius: 6px; transition: all 0.2s; }
  .btn-remove { background: #fee2e2; color: #ef4444; height: 38px; width: 38px; display: flex; align-items: center; justify-content: center; }
  .btn-remove:hover { background: #fca5a5; }
  .form-actions { display: flex; justify-content: space-between; margin-top: 20px; padding-top: 16px; border-top: 1px solid #f1f5f9; }
  .btn-add { background: #f1f5f9; color: #475569; padding: 10px 16px; }
  .btn-add:hover { background: #e2e8f0; }
  .btn-submit { background: #3b82f6; color: white; padding: 10px 24px; font-size: 1rem; }
  .btn-submit:hover:not(:disabled) { background: #2563eb; }
  .btn-submit:disabled { opacity: 0.7; cursor: not-allowed; }
</style>