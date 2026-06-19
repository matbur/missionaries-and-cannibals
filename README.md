# missionaries-and-cannibals

A solver for the classic **Missionaries and Cannibals** river-crossing puzzle, implemented in Go.

## Rules

Three missionaries (**M**) and three cannibals (**C**) start on the **left bank** of a river. A boat that holds at most two people must ferry everyone to the **right bank**.

### Goal

Move all six people to the right bank:

```
Start:  Left: 3M 3C  |  Right: 0M 0C  |  Boat: left
Goal:   Left: 0M 0C  |  Right: 3M 3C  |  Boat: right
```

### Boat

- The boat carries **1 or 2 people** per trip (never empty).
- Allowed moves: 2M, 1M, 1M+1C, 1C, 2C.
- The boat travels back and forth between the two banks.

### Safety

On **both banks**, missionaries must not be outnumbered by cannibals (unless there are no missionaries on that bank).

If this rule is broken, the path is invalid:

| Outcome       | Meaning                                              |
|---------------|------------------------------------------------------|
| `EATEN_LEFT`  | Cannibals outnumber missionaries on the **right** bank |
| `EATEN_RIGHT` | Cannibals outnumber missionaries on the **left** bank  |

### Other constraints

- No bank may have more than 3 people of either type.
- Revisiting the same state is forbidden (`LOOP`).

## How it works

The program explores all valid move sequences using breadth-first search over a tree of states. It prints a summary of explored paths and lists every distinct winning solution.

## Usage

```
make run
```

Build a binary:

```
make build
```

Regenerate protobuf code:

```
make proto
```

Other targets: `make help`, `make tidy`, `make clean`.
