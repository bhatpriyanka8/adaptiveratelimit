# Design Principles

This document describes the design philosophy and constraints behind
`adaptiveratelimit`.

The goal of this library is to provide **simple, safe, and predictable
adaptive rate limiting for a single process**.

---

## Core Principles

### Bounded Behavior
- Rate limits must always respect configured minimums and maximums.
- The limiter must not grow unbounded state over time.

### Adaptive but Stable
- Adaptation should respond to sustained trends, not transient spikes.
- EWMA is used to smooth latency and error signals.
- Cooldowns are required to prevent oscillation and thrashing.

### Fast Hot Path
- The request path (`Allow`, `Record`) should remain lightweight.
- Expensive computation belongs in background control loops.

### Concurrency Safety
- All public APIs must be safe for concurrent use.
- Background goroutines must be owned and managed by the limiter.

### Explicit Backpressure
- When the limit is exceeded, requests are explicitly rejected.
- Silent dropping or implicit blocking is discouraged.

### Minimal API Surface
- Public APIs should be small and intentional.
- Internals should not be exposed unless there is a clear use case.

### Clean Shutdown
- Background loops must stop when `Stop()` is called.
- No goroutine leaks or runaway timers are acceptable.

---

## Non-Goals

This project intentionally does **not** aim to provide:

- Distributed or global rate limiting
- Cross-process coordination
- Fairness across multiple clients
- Persistent state across restarts

These are better handled by centralized or distributed systems.

---

## Trade-offs

- Instance-local design favors simplicity and predictability over global fairness.
- EWMA smoothing favors stability over immediate reaction.
- Cooldowns trade faster adaptation for safer behavior.

These trade-offs are intentional.
