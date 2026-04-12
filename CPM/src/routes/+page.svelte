<script>
  import Graph from '$lib/components/Graph.svelte';

  let activities = $state([]);
  let loading = false;

  async function fetchActivities() {
    loading = true;
    try {
      const response = await fetch('http://localhost:8080/activity', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify([
    {"id": "A", "duration": 3, "predecessors": []},
    {"id": "B", "duration": 5, "predecessors": ["A"]},
    {"id": "C", "duration": 2, "predecessors": ["A"]},
    {"id": "D", "duration": 4, "predecessors": ["B"]},
    {"id": "E", "duration": 6, "predecessors": ["C"]},
    {"id": "F", "duration": 3, "predecessors": ["B"]},
    {"id": "G", "duration": 5, "predecessors": ["D", "E"]},
    {"id": "H", "duration": 2, "predecessors": ["F"]},
    {"id": "I", "duration": 4, "predecessors": ["G", "H"]},
    {"id": "J", "duration": 3, "predecessors": ["I"]},
    {"id": "K", "duration": 2, "predecessors": ["J"]}
        ])
      });
      activities = await response.json();
    } catch (e) {
      console.error("Błąd pobierania danych", e);
    } finally {
      loading = false;
    }
  }
</script>

<button onclick={fetchActivities}>Oblicz ścieżkę</button>

<Graph {activities} />

<pre>{JSON.stringify(activities, null, 2)}</pre>
