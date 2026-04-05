import requests
import json

url = "http://localhost:8080/event"

tasks = [
    {
        "id": "A",
        "from_node": 1,
        "to_node": 2,
        "duration": 3
    },
    {
        "id": "B",
        "from_node": 1,
        "to_node": 3,
        "duration": 5
    },
    {
        "id": "C",
        "from_node": 2,
        "to_node": 4,
        "duration": 4
    },
    {
        "id": "D",
        "from_node": 3,
        "to_node": 4,
        "duration": 6
    }
]

print("\n- - - - - ZDARZENIOWY - - - - -\n")

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

print ("\n- - - - - - - - - - - - - - -\n")