# Roadmap

This document outlines the planned and potential future direction of
`adaptiveratelimit`.

The roadmap is **directional, not a commitment**. Features may evolve,
be deferred, or be dropped based on real-world usage and maintainability.

---

## v0.2.x — Observability & Usability

This release focuses on making the limiter easier to reason about and
operate, without changing its core behavior.

### Planned

- Additional read-only observability APIs
  - Expose current adjustment direction (increase / decrease / steady)
  - Expose time since last limit change

- Configuration validation
  - Detect invalid or unsafe configurations early
  - Fail fast on startup instead of degrading at runtime

- Improved documentation and examples
  - More examples
  - Expanded design notes and usage guidance

- Minor API polish
  - Clarify method names and comments where needed
  - Keep backward compatibility

---

## v0.3.x — Extensibility & Control

This release focuses on making the adaptive logic more flexible while
preserving safe defaults.

### Planned

- Pluggable adaptation strategy
  - Allow custom controllers to be injected
  - Keep EWMA-based strategy as the default

- Adjustable control loop behavior
  - Separate sampling interval from adjustment interval
  - More control over cooldown behavior

- Optional context-aware APIs
  - Support `context.Context` for request lifecycle awareness
  - Enable better integration with cancellations and deadlines

---

## Future / Exploratory (No Timeline)

These ideas are under consideration but are **not scheduled**.

- Adaptive burst handling
  - Token or credit-based bursts with bounded limits

- Percentile-based adaptation
  - Adjust limits based on p95 / p99 latency instead of averages

- Optional metrics hooks
  - User-provided hooks for exporting metrics (Prometheus, etc.)
  - Keep metrics opt-in and dependency-free

- Request classification
  - Different limits or behavior based on request type or priority

- Instance-level coordination
  - Coordinating limits across goroutines within a process
  - No cross-process or distributed coordination

---

## Non-Goals

The following are intentionally out of scope:

- Distributed or global rate limiting
- Cross-node coordination
- Persistent state across restarts
- Hard real-time guarantees

These problems are better solved by dedicated distributed systems.