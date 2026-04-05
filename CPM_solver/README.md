# CPM solver
Napisany w Golang serwis. Przyjmuje POST na `:8080` i oddaje json z rozwiązaniem metody.

1. [Uruchamianie](#uruchamianie)
2. [Legenda](#legenda)
3. [Wariant czynnościowy](#wariant-czynnościowy)
   - [Przykład](#przykład)
   - [Test](#test)
4. [Wariant zdarzeniowy](#wariant-zdarzeniowy)
   - [Przykład](#przykład-1)
   - [Test](#test-1)

# Uruchamianie 
Skompilowany `.exe` wystarczy odpalić. Można też `go mod tidy`, po czym `go run main.go`.

# Legenda
- ES (Early Start) - najwcześniejszy moment rozpoczęcia zadania.
- EF (Early Finish) - najwcześniejszy moment zakończenia.
- LS (Late Start) - najpóźniejszy bezpieczny moment rozpoczęcia.
- LF (Late Finish) - najpóźniejszy bezpieczny moment zakończenia.
- Slack - rezerwa czasu (luz). Jeśli wynosi 0, czynność jest krytyczna.
- ET (Earliest Time) - najwcześniejszy czas wystąpienia zdarzenia (węzła).
- LT (Latest Time) - najpóźniejszy dopuszczalny czas wystąpienia zdarzenia.

# Wariant czynnościowy
Endpoint: `/activity`.

## Przykład
Przykładowy Post z takim json w Body:
```json
[
    {"id": "A", "duration": 5, "predecessors": []},
    {"id": "B", "duration": 3, "predecessors": ["A"]},
    {"id": "C", "duration": 2, "predecessors": ["A"]},
    {"id": "D", "duration": 4, "predecessors": ["B", "C"]}
]
```

Przykładowa odpowiedź
```json
[
    {
        "id": "A",
        "duration": 5,
        "predecessors": [],
        "es": 0,
        "ef": 5,
        "ls": 0,
        "lf": 5,
        "slack": 0,
        "is_critical": true
    },
    {
        "id": "B",
        "duration": 3,
        "predecessors": [
            "A"
        ],
        "es": 5,
        "ef": 8,
        "ls": 5,
        "lf": 8,
        "slack": 0,
        "is_critical": true
    },
    {
        "id": "C",
        "duration": 2,
        "predecessors": [
            "A"
        ],
        "es": 5,
        "ef": 7,
        "ls": 6,
        "lf": 8,
        "slack": 1,
        "is_critical": false
    },
    {
        "id": "D",
        "duration": 4,
        "predecessors": [
            "B",
            "C"
        ],
        "es": 8,
        "ef": 12,
        "ls": 8,
        "lf": 12,
        "slack": 0,
        "is_critical": true
    }
]
```

## Test
W folderze test znajduje się pythonowy unit test z powyższym casem. Plik: `test_activity.py`.

# Wariant zdarzeniowy
Endpoint: `/event`.

## Przykład
Przykładowy Post z takim json w Body:
```json
[
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
```

Przykładowa odpowiedź
```json
{
    "events": [
        {
            "node": 1,
            "et": 0,
            "lt": 0,
            "slack": 0,
            "is_critical": true
        },
        {
            "node": 2,
            "et": 3,
            "lt": 7,
            "slack": 4,
            "is_critical": false
        },
        {
            "node": 3,
            "et": 5,
            "lt": 5,
            "slack": 0,
            "is_critical": true
        },
        {
            "node": 4,
            "et": 11,
            "lt": 11,
            "slack": 0,
            "is_critical": true
        }
    ],
    "activities": [
        {
            "id": "A",
            "from_node": 1,
            "to_node": 2,
            "duration": 3,
            "es": 0,
            "ef": 3,
            "ls": 4,
            "lf": 7,
            "slack": 4,
            "is_critical": false
        },
        {
            "id": "B",
            "from_node": 1,
            "to_node": 3,
            "duration": 5,
            "es": 0,
            "ef": 5,
            "ls": 0,
            "lf": 5,
            "slack": 0,
            "is_critical": true
        },
        {
            "id": "C",
            "from_node": 2,
            "to_node": 4,
            "duration": 4,
            "es": 3,
            "ef": 7,
            "ls": 7,
            "lf": 11,
            "slack": 4,
            "is_critical": false
        },
        {
            "id": "D",
            "from_node": 3,
            "to_node": 4,
            "duration": 6,
            "es": 5,
            "ef": 11,
            "ls": 5,
            "lf": 11,
            "slack": 0,
            "is_critical": true
        }
    ]
}
```

## Test
W folderze test znajduje się pythonowy unit test z powyższym casem. Plik: `test_event.py`.
