# TODO

### Branch and PR plan
- Branch name: `feat/streaming-ingestion-backfill`
- PR sequence:
  1) Scaffolding, config, interfaces, storage substrate
  2) WebSocket streaming (accountSubscribe, logsSubscribe) + connection manager
  3) Backfill worker (getSignaturesForAddress) + durable cursors
  4) Deduplication + resume logic integration
  5) End-to-end integration, metrics, docs, benchmarks

### Proposed structure (new packages/files)
- `internal/streaming/`
  - `ws_manager.go`: pooled WS client, backoff, resubscribe, health
  - `subscriptions.go`: accountSubscribe/logsSubscribe helpers
  - `events.go`: `Event` model (account/log/tx), decode utilities
- `internal/backfill/`
  - `worker.go`: cursor-based `getSignaturesForAddress` walker
  - `pager.go`: pagination windowing and rate-limit handling
- `internal/dedupe/`
  - `signature_store.go`: durable signature set with in-memory LRU+persisted
- `internal/cursors/`
  - `cursor_store.go`: durable cursors per address/program
- `internal/storage/`
  - `kv.go`: simple KV store (bbolt or badger) abstraction
- `cmd/monitor/`
  - integrate new streaming/backfill runners behind config flags
- Update: `config.example.json`, `README.md`, `docs/roadmap.json`

### Config additions
- In `internal/config/config.go`:
  - `Streaming`: `enabled`, `rpc_ws_url` (default derive from `network_url`), `program_ids`, `wallets`, `pool_size`, `max_retries`, `backoff`
  - `Backfill`: `enabled`, `start_time`, `end_time` or `since_signature`, `batch_size`, `concurrency`
  - `Dedup`: `enabled`, `retention_hours`, `max_in_memory`, `persist_path`
  - `Resume`: `gap_fill_on_reconnect` (bool), `max_gap_seconds`
- Update `config.example.json` with sane defaults.

### Tasks by PR

#### PR 1: Scaffolding, config, interfaces, storage substrate
- Add config structs and validation.
- Introduce `internal/storage/kv.go` with an interface:
  - `Put(key []byte, value []byte) error`, `Get(key []byte) ([]byte, error)`, `Delete`, `Iterate(prefix []byte, f func(k,v []byte) bool)`
- Implement KV using `bbolt` (pure-Go, durable) with a single DB file `./data/state.db`.
- Define interfaces:
  - `Event` (type, signature, slot, blockTime, account, programId, payload)
  - `StreamSource` (Start/Stop, Subscribe, Health)
  - `BackfillSource` (RunRange, RunSinceCursor)
  - `CursorStore` (Get/Set per address/program)
  - `SignatureStore` (Has/Add, TTL pruning)
- Wire feature flags (streaming/backfill/dedup/resume) from config.
- Minimal docs and diagrams.

Acceptance alignment: sets foundations for real-time, backfill, dedupe, resume.

#### PR 2: WebSocket streaming ingestion
- Add WS manager using `github.com/gagliardetto/solana-go/rpc/ws`.
  - Pooled connections from configurable endpoints (fallback to main if only one).
  - Exponential backoff with jitter; health pings; metrics counters.
- Implement `accountSubscribe` for configured wallets.
  - Decode account changes; produce `Event` with wallet, mint, lamports/token changes.
- Implement `logsSubscribe` for configured `program_ids`.
  - Parse logs, extract signatures; enrich by fetching tx meta only when needed.
- Event router → channel to processing pipeline (later shared with backfill).
- Basic latency measurement (recv time – on-chain block time if available) to log.

Acceptance alignment: WebSocket streaming operational; aim for <1s latency.

#### PR 3: Backfill worker with durable cursors
- Implement `getSignaturesForAddress` walker:
  - Window by time and by page (before/limit), configurable `batch_size`.
  - For each signature, fetch confirmed data when needed to build `Event`.
- `CursorStore`: store last processed signature per address/program with block time.
- Run modes:
  - Range mode: `start_time`→`end_time`.
  - Catch-up mode: from `cursor` to now; then signal “caught-up” to hand off to streaming steady-state.
- Rate limit handling with exponential backoff and concurrency control.

Acceptance alignment: Historical backfill for any time range; performance goal path to >1000 tx/s via concurrency and minimal per-tx RPC.

#### PR 4: Deduplication and resume logic
- `SignatureStore`: in-memory LRU (e.g., map+ring or ristretto) plus persisted set in KV.
  - `Has(signature)`, `Add(signature)`, background compaction by `retention_hours`.
- Processor wrapper ensuring exactly-once per signature across streaming and backfill.
- Resume:
  - On WS disconnect, mark gap start; on reconnect, compute time/slot gap, run targeted backfill from last cursor to now; then resubscribe.
  - Ensure no duplicate events during gap fill.
- Add smoke tests for duplicate prevention and reconnect flows.

Acceptance alignment: Zero duplicates; graceful reconnection with gap fill.

#### PR 5: Integration, monitoring, docs, benchmarks
- Integrate with existing `cmd/monitor/main.go`:
  - Start backfill (if configured) → wait until caught up → start streaming; or run both with dedupe if desired.
  - Optionally keep periodic scans as a fallback-only feature; otherwise, rely on events to drive alerting.
- Logging: use existing `utils.Logger`; add counters:
  - `events_received`, `events_processed`, `duplicates_dropped`, `reconnects`, `gapfill_events`, `avg_latency_ms`, `backfill_tps`.
- Add CLI flags to override config for quick testing.
- Benchmarks:
  - Synthetic replay mode to measure throughput
  - Collect latency percentiles and backfill TPS; document results in `docs/insider-monitor-improvement-report.md`.
- Update `README.md` and `config.example.json`.

Acceptance alignment: Comprehensive logging/monitoring; performance benchmarks documented.

### Testing plan
- Unit tests:
  - Cursor store semantics (idempotent Set/Get, per address).
  - Signature store dedupe across restarts.
  - WS manager reconnection and resubscribe logic (mock ws).
  - Backfill pagination (edges: empty page, duplication when crossing time windows).
- Integration tests (devnet):
  - Stream an address with known activity; verify latency <1s typical.
  - Backfill a known range; verify dedupe=0 duplicates; measure TPS.
- Fault injection:
  - Force WS disconnects; confirm gap fill and no data loss.

### Observability and performance
- Emit periodic logs with moving averages:
  - `events/sec (stream)`, `events/sec (backfill)`, `latency_ms p50/p95`, `duplicates_dropped`, `reconnects`.
- Optional: Prometheus endpoint behind a flag in a future PR.

### Risks and mitigations
- Rate limits: use backoff, batching, and optional provider guidance already in `config`.
- WS instability: pool + auto-resubscribe + gap fill.
- Disk contention: `bbolt` single-writer; mitigate by batching writes and using a fan-in writer goroutine.
- Throughput: avoid per-event heavy RPC; fetch tx data only when necessary; parallel fetch with bounded concurrency.

### Definition of done (mapped to Acceptance Criteria)
- WebSocket streaming: verified with observed end-to-end latency <1s median.
- Backfill: configurable time range; measured TPS ≥ 1000 on devnet with batching.
- Dedup: no duplicates across backfill + streaming + reconnect tests.
- Resume: forced disconnects recover with zero data loss; cursors updated.
- Logs/metrics: counters present; benchmarks included in docs.

### Minimal code touchpoints in current repo
- Keep `internal/monitor` intact; add a new event-driven path. Over time, alerts can subscribe to the event bus rather than relying on periodic scans.
- Extend `internal/config` and `config.example.json`.
- Add `internal/storage/kv.go`, `internal/cursors`, `internal/dedupe`, `internal/streaming`, `internal/backfill`.
- Wire new components in `cmd/monitor/main.go` behind config flags.

### Next steps (quick start)
- Create branch.
- Land PR 1 with interfaces, KV, config.
- Land PR 2 streaming on devnet with 1-2 wallets and 1 program id.
- Land PR 3 backfill with cursors; confirm TPS targets.
- Land PR 4 dedupe + resume; prove zero dupes under churn.
- Land PR 5 integration + benchmarks + docs.

- I’ve identified the current polling flow and where the streaming/backfill layers will integrate. Next I’d scaffold interfaces and config in a first PR, then proceed to streaming, backfill, and dedupe/resume in sequence.

- Key changes will add `internal/streaming`, `internal/backfill`, `internal/dedupe`, `internal/cursors`, and a KV store in `internal/storage`. We’ll introduce config for WS, backfill, dedupe, and resume, integrate them in `cmd/monitor/main.go`, and provide metrics and docs meeting the acceptance criteria.