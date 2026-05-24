# Hexa v1.0.0 — Release Plan & Engineering Notes

> Working/handoff document. Written so a fresh session (with no memory of prior
> work) can pick up taking `github.com/kamva/hexa` to a stable v1.0.0.
> Last updated during the hardening session described below.

---

## 0. Orientation for a new session

- **Repo:** `github.com/kamva/hexa` (GitHub: `Kamva/hexa`). A hexagonal-architecture
  **microservice SDK** for the Kamva ecosystem (sibling repos: `hexa-event`,
  `hexa-arranger`/Temporal, `hexa-rpc`, `hexa-tuner`, `hexa-sendo`, `hexa-job`,
  `hexa-echo`).
- **Default branch is `master`** (NOT `main`) — confirmed via
  `git ls-remote --symref origin HEAD`. PRs target `master`.
- **GitHub access:** use the `mcp__github__*` tools (owner `Kamva`, repo `hexa`).
  There is **no `gh` CLI** and no direct GitHub API in this environment.
- **It's pre-1.0** and self-describes as "use at your own risk" in the README.
- This document is the plan to make it trustworthy enough to tag **v1.0.0**.

### Fast local verification (matches CI)
```bash
gofmt -l .                 # must be empty (CI checks the WHOLE tree, incl. examples/)
go vet ./...
go build ./...
go test ./... -race -covermode=atomic -coverprofile=cover.out
go tool cover -func=cover.out | tail -1   # total coverage
# integration tests skip unless these are set:
HEXA_TEST_MONGO_URI=mongodb://...  go test ./hdlm/mongolock/...
HEXA_TEST_REDIS_ADDR=127.0.0.1:6379 go test ./hdlm/redislock/...
```

---

## 1. State of the repo (after the hardening session)

A code review found ~16 issues; the genuine bugs were fixed, CI was added, and
tests were backfilled. **All of the following PRs are merged into `master`.**

### Bug fixes (merged)
| PR | Area | Fix |
|----|------|-----|
| #6 | `hurl` | Repaired non-compiling package (`io` import, `HTTPErr`, exported `ResponseErr`, close body); `NewClient("")` no longer panics on relative paths |
| #7 | `sr` | Multi-registry off-by-one (`i > 0` → `i >= 0`); stop clobbering an explicitly-set `Descriptor.Health` |
| #8 | `hdlm/redislock` | `Unlock` now actually releases (was a no-op); keeps lock handle, refresh-on-relock, owner-as-token; init embedded `Health`; preserve handle on transient `Release` error |
| #9 | `hlog` | `SetGlobalLogger` binds `Message` to driver `Message` (was `Info`); `FieldToKeyVal` decodes real Go values via `zapcore.NewMapObjectEncoder` |
| #10 | `lg` | `UseEmbeddedFieldsInPackage` returns the repackaged slice (was returning the original) |
| #11 | `hexa` (root) | `WithBaseTranslator` writes the correct context key; context propagation no longer requires a user; `validateUserMetaData` checks string-field types |
| #12 | `hexa` | `defaultError.Unwrap()` so `errors.Is/As` traverse the internal error |
| #13 | `hurl` | Redact `Authorization`/`Proxy-Authorization`/`Cookie`/`Set-Cookie` in request/response debug dumps |
| #5 (earlier) | `hlog/logdriver` | Sentry captures the original error (preserves stack/wrap chain) |

### Infra / docs (merged)
| PR | What |
|----|------|
| #14 | `CHANGELOG.md` (Keep a Changelog) + README min Go version `1.13` → `1.18` |
| #15 | **First CI** (`.github/workflows/ci.yml`) + mongo test skip-guard |

### Test backfill (merged) — coverage before → after
| PR | Package | Coverage |
|----|---------|----------|
| #16 | `sr` | 0% → 75.3% |
| #17 | `hurl` | 1.8% → 72.1% |
| #18 | `probe` | 0% → 71.4% |
| #19 | `hexatranslator` | 0% → 72.7% |
| #20 | `hlog` | 5.2% → 92.2% |
| #21 | `lg` | 8% → 52.7% |
| #22 | `hexa` (root) | 44.3% → 59.1% |
| #23 | `htel` | 27.8% → 100% |
| #24 | `db/mgmadapter` | 13% → 54.3% |
| #25 | `hdlm/redislock` | 0% → 18.4% (no redis; higher with `HEXA_TEST_REDIS_ADDR`) |

Library total was ~22%; now much higher. Several PRs include regression tests
pinning the bugs above.

### CI shape (current)
- Job **`build & test`** (enforced gate): gofmt check → `go vet` → `go build` →
  `go test ./... -race -covermode=atomic -coverprofile` → coverage summary.
  Runs on **Go 1.22** (pinned; see gotcha below).
- Job **`lint (advisory)`**: `golangci-lint` pinned to **v1.55.2** on **Go 1.21**,
  run with `--issues-exit-code=0` so it's **green/logs-only** (does not block).
- Triggers: `push` to `master` and all `pull_request`s.

---

## 2. Environment gotchas (these cost real time — read before touching CI/lint)

1. **golangci-lint version vs config:** `.golangci.yaml` is the **v1 format**
   (header comment says "Expected golangci-lint version: 1.49.0"). The
   golangci-lint binary installed in this sandbox is **v2** and refuses the v1
   config ("unsupported version"). CI therefore pins **golangci-lint v1.55.2**.
2. **golangci-lint v1.55.2 cannot run under Go 1.24** locally — it emits spurious
   `typecheck` errors ("could not import ... unsupported version: 2") because it
   can't read Go 1.24 export data. It needs **Go ≤ 1.21**. CI's lint job uses
   `actions/setup-go` with `go-version: "1.21"` for this reason.
3. **Cannot download old Go SDKs here:** `go install golang.org/dl/go1.21.x` then
   `goX download` **fails** (dl.google.com returns 403). So lint cannot be fully
   verified locally in this sandbox; rely on CI for the real lint result.
4. **`gofmt` whole-tree:** CI's formatting step runs `gofmt -l .` over everything,
   **including `examples/`**. `examples/telemetry/main.go` was unformatted and
   broke the first CI run (now fixed). Always `gofmt -w` new/changed files.
5. **Go 1.18 vs deps:** `go.mod` declares `go 1.18`, but CI's test job is pinned to
   **1.22** because it's unverified that the dependency graph builds on 1.18 (the
   sandbox only has Go 1.24). If you want to honor the declared minimum, test a
   `[1.18, 1.22]` matrix in CI and see whether 1.18 actually builds.
6. **WebFetch is blocked (403)** for `go-kratos.dev`, `gokit.io`,
   `uber-go.github.io`. Use `WebSearch` summaries instead (see §5).
7. **Integration tests** (`hdlm/mongolock`, `hdlm/redislock`) **skip** unless
   `HEXA_TEST_MONGO_URI` / `HEXA_TEST_REDIS_ADDR` are set. Wire these to CI
   service containers to actually exercise the lock drivers.
8. **Git push:** use `git push -u origin <branch>`; retry on network errors with
   backoff. Don't push to `master` directly. Each change so far went via a PR.

---

## 3. Known issues still OPEN (the v1 punch list)

Deferred during the bug-fix pass — none are build-breaking, but they matter for a
stable release:

- **`sr` registry concurrency:** `serviceRegistry.l` (the descriptor slice) is
  mutated by `Register*` and read by `Boot`/`Service`/`Descriptors`/shutdown with
  no mutex (only `booted`/`done` are atomic). Add a mutex if concurrent
  registration is to be supported. (`hexa.Store` is already concurrency-safe;
  this is inconsistent.)
- **`hlog.LevelFromString` panics** on an unknown string — should return an error
  or a default for config parsing.
- **`hurl.ResponseErr` uses `r.StatusCode <= 300`** — treats HTTP 300 as success
  and 3xx>300 as errors. Decide the intended threshold (likely `< 400`).
- **`go.mod` hygiene:** `github.com/bsm/redislock` and `github.com/redis/go-redis/v9`
  are listed `// indirect` though directly imported by `hdlm/redislock`. Run
  `go mod tidy` (NOT done yet — was deliberately skipped to avoid toolchain churn).
- **Version mismatch:** `version.go` declares `Version = "1.1.0"` but the only git
  tag is `v0.1.0`. Reconcile before tagging v1.
- **README vs CI Go version:** README now says min Go `1.18`; CI tests on `1.22`.
  Pick one story (supported floor vs tested versions) and make them consistent.
- **`session.go` is marked "prototype"** — finish, gate behind an `experimental`
  package, or exclude from the v1 API surface.
- **Sentry driver `setUser`** calls `gutil.IP(r)` where `r` may be nil (user
  present, request absent) — possible nil deref; verify/guard.
- **Untested code:** `hlog/logdriver` zap driver; `db/mgmadapter` `adapter.go`
  (index creation) + `monitoring.go` (need Mongo); `examples/` (build-only).
- **Lint backlog:** most packages have **no package comment** while `revive`'s
  `package-comments` rule is enabled — that's the main reason lint is advisory.
  Add package comments to flip lint to a hard gate (see Phase 4).

---

## 4. The v1.0.0 plan

**Definition of v1:** a documented, tested, dependency-current release whose
exported API is **frozen under semver**. Breaking changes after v1 require the
`github.com/kamva/hexa/v2` module path (Go modules rule). v1 is about
*trustworthiness*, not new features.

### Phase 0 — Foundation ✅ (done this session)
CI gate, test backfill, bug fixes, CHANGELOG. Remaining cleanup:
- [ ] Clear the package-comment backlog, then **flip lint to blocking**
      (remove `--issues-exit-code=0` and set the job as a required check).
- [ ] Add Mongo/Redis **service containers** to a CI integration job and set
      `HEXA_TEST_MONGO_URI` / `HEXA_TEST_REDIS_ADDR`.

### Phase 1 — API surface & stability policy
- [ ] Audit every exported symbol; label each package **stable** or **experimental**.
- [ ] `COMPATIBILITY.md` / stability section: what v1 guarantees, semver rules,
      `/v2` for future breaks (Kratos model).
- [ ] Reconcile `version.go` (`1.1.0`) vs tag (`v0.1.0`); decide whether to keep a
      version constant or derive from build info.
- [ ] `go mod tidy`; declare a supported **Go floor** (suggest 1.21+); reflect in
      `go.mod`, README, CI matrix.
- [ ] Decide fate of `session.go` (prototype) and any rough `lg` surface.

### Phase 2 — Correctness & hardening
- [ ] Close coverage gaps: zap log driver, `mgmadapter` Mongo paths (integration),
      `examples` build-smoke.
- [ ] Land deferred items from §3: registry mutex, `LevelFromString` error,
      `hurl` status threshold, `gutil.IP(nil)` guard.
- [ ] **Modernize dependencies** (riskiest workstream; do early in an RC):
      otel `v1.2.0` (2021) → current, `mongo-driver 1.7.0` → current, `zap 1.14`,
      `go-i18n 2.0.3`. Re-verify under new versions. Note: bumping otel will touch
      `htel/` (propagator + telemetry wrappers).

### Phase 3 — API polish (informed by peer frameworks, see §5)
- [ ] **Errors:** consider Kratos-style split — a stable machine `Reason`/ID
      distinct from a transport `Code` — and act on the repo's own README proposal
      to map **gRPC status codes**. If breaking, do it **now (pre-1.0)**.
      Current `hexa.Error` = `httpStatus` + `id` + `localizedMessage` + `data` +
      `reportData`; `errors.Is/As` now works via `Unwrap` (#12).
- [ ] **Observability:** ensure `htel` integrates with current otel; verify hexa
      context propagation ↔ otel span context.
- [ ] Consistency: package comments, receiver naming, globals concurrency
      (`hlog.SetGlobalLogger`, `hexatranslator.SetGlobal` mutate package vars
      without sync — fine at startup, racy if used later).

### Phase 4 — Documentation & DX (the biggest gap vs Kratos)
- [ ] Package-level godoc on every package (also unblocks lint `package-comments`).
- [ ] `docs/`: a **Design/Philosophy** page (Kratos-style) + per-concept guides:
      Context & propagation, Errors/Reply, Health/Probe, Service Registry, DLM,
      Logging (`hlog`), Translator, `hurl`, `lg` (layer generator) + a **Quickstart**.
- [ ] Real, building `examples/` per subsystem (CI-built); a `hexa-tuner`
      "boot everything" walkthrough.
- [ ] README overhaul: CI/coverage/go-reference badges, correct Go version,
      concept map, links to sibling repos.

### Phase 5 — Release mechanics
- [ ] Cut **`v1.0.0-rc.1`**; dogfood in sibling repos (`hexa-echo`, `hexa-event`,
      `hexa-tuner`) to validate the frozen API against real consumers.
- [ ] Coordinate the ecosystem bump (the `hexa-*` modules depend on `hexa`).
- [ ] Tag **`v1.0.0`**; GitHub release from CHANGELOG; add `CONTRIBUTING.md` and
      release automation (goreleaser or a tag workflow). Dependabot already present.

**Sequencing:** Phases 1–2 can run in parallel after lint/integration cleanup;
Phase 4 (docs) is the long pole — start it early. Target **rc.1 after Phases 1–3**,
then docs + dogfooding, then **1.0.0**.

### Open decisions for the maintainer
1. **API freeze scope** — is `session` (and rough `lg`) in or out of v1?
2. **Errors** — adopt gRPC-code mapping / stable-Reason split now (breaking,
   cleaner v1) or defer to v2?
3. **Ecosystem** — release `hexa` v1 standalone or in lockstep with `hexa-*`?
4. **Go floor** — 1.21, 1.22, or a 1.18→latest matrix?

---

## 5. Peer-framework research notes (concepts to borrow)

Reading blocked by 403s; these are from web-search summaries (sources below).

- **Kratos** (`go-kratos`, closest analogue): layered design (lifecycle /
  transport / cross-cutting / codegen / contrib); "socket" philosophy —
  standardized, pluggable interfaces, not bound to any infra. **Errors are a
  first-class model**: `Code` (HTTP-like int32) + `Reason` (stable machine string)
  + `Message` (human) + `Metadata` (map), generated from **proto**. Observability
  (trace/metrics/log) built in. Breaking changes live behind **`/v2`** module path.
  → Hexa shares the pluggable ethos; lacks concept docs, the stable-reason vs
  transport-code split, and a versioning policy.
- **go-kit:** strict **service / endpoint / transport** separation; business logic
  knows nothing about transports; middleware at endpoint + transport layers.
  → Validates Hexa's interface-first, transport-agnostic instincts.
- **go-micro:** "**sane defaults, pluggable architecture**," one Go interface per
  distributed-system abstraction (registry / transport / broker).
  → Same ethos as Hexa's DLM / registry / driver pattern.
- **Uber Fx:** DI via `Provide`/`Invoke` + a `Lifecycle` of ordered
  `OnStart`/`OnStop` hooks run in dependency order with timeouts.
  → More sophisticated than Hexa's priority-ordered `ServiceRegistry`. Hexa
  deliberately punts DI to `wire`; make that boundary explicit in v1 docs.
- **Go modules versioning:** v0/v1 stay at the bare module path; first stable tag
  is `v1.0.0`; any later breaking change needs `/v2`. → v1 is fundamentally a
  promise about the exported API.

### Sources
- Kratos: https://go-kratos.dev/ , https://go-kratos.dev/docs/intro/design/ ,
  https://go-kratos.dev/docs/component/errors/
- go-kit: https://gokit.io/faq/
- go-micro: https://github.com/micro/go-micro
- Uber Fx: https://uber-go.github.io/fx/
- Go modules (v2+ paths/semver): https://go.dev/wiki/Modules

---

## 6. Architecture cheat-sheet (what lives where)

- **Root `hexa`** package: `Context` (user/correlation-id/locale/logger/translator/
  store, auto-rederives logger & translator on change), `ContextPropagator`
  (default + multi + keys), `Error`/`Reply`, `User`/`UserPropagator`, `Health`/
  `HealthReporter`, `ServiceRegistry` interfaces, `DLM`/`Mutex`, `Store`, `ID`,
  `Secret`, `Translator` interface, `Runnable`/`Bootable`/`Shutdownable`.
- **`sr`**: `ServiceRegistry` impl (priority-ordered boot, reverse-order shutdown)
  + `multiSearchRegistry`.
- **`hlog`**: logging facade + drivers in `hlog/logdriver` (zap, sentry, stack);
  printer driver for tests; package-level global logger funcs.
- **`hexatranslator`**: `Translator` impls (i18n via go-i18n, empty, key) + global.
- **`htel`**: OpenTelemetry wrappers + hexa↔otel propagator (`HexaCarrier`).
- **`hurl`**: improved `net/http` client (base URL, options, logging w/ redaction).
- **`hdlm`**: distributed lock managers — `mongolock` (clever upsert-or-dup-key
  query) and `redislock` (bsm/redislock).
- **`db/mgmadapter`**: MongoDB (mgm) adapter — `hexaID`, `IDField`, `Touchable`,
  index/monitoring helpers.
- **`pagination`**, **`hexamask`** (PATCH field masks), **`probe`** (health/pprof
  HTTP server), **`lg`** (layer/code generator over go/ast).
- Notable: the `Runnable` interface doc explains its non-blocking redesign;
  `mongolock` `TryLock` filter `{_id, $or:[{owner:me},{expiry<now}]}` + upsert is
  the correct contention-detection trick.
