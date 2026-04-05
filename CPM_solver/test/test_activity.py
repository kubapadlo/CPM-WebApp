import requests
import json

url = "http://localhost:8080/activity"

tasks = [
    {"id": "A", "duration": 5, "predecessors": []},
    {"id": "B", "duration": 3, "predecessors": ["A"]},
    {"id": "C", "duration": 2, "predecessors": ["A"]},
    {"id": "D", "duration": 4, "predecessors": ["B", "C"]}
]

print("\n- - - - - CZYNNOSCIOWY - - - - -\n")

try:
    response = requests.post(url, json=tasks)

    if response.status_code == 200:
        print("Sukces! Wyniki obliczeń:")
        print(json.dumps(response.json(), indent=4))
    else:
        print(f"Błąd! Kod statusu: {response.status_code}")
        print("Odpowiedź serwera:", response.text)

except Exception as e:
    print(f" Nie udało się połączyć z serwerem: {e}")

print("\n- - - - - - - - - - - - - - - - -\n")