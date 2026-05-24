# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
This project is pre-1.0 and does not yet follow semantic versioning.

## [Unreleased]

### Security

- **hurl:** Sensitive headers (`Authorization`, `Proxy-Authorization`, `Cookie`,
  `Set-Cookie`) are now redacted as `REDACTED` when requests/responses are dumped
  to logs via the `LogMode*Headers` flags, so credentials and session data no
  longer leak into logs. (#13)

### Fixed

- **hurl:** Repaired a broken build in `error.go` — `io` was used without being
  imported, the error literal referenced an undefined `HttpErr` (the type is
  `HTTPErr`), and the function was named `responseErr` while its only caller used
  the exported `ResponseErr`. The error path now also closes the response body it
  reads. (#6)
- **hurl:** `NewClient` created with an empty base no longer panics on a relative
  request path; it returns a clear error instead. (#6)
- **sr:** `multiSearchRegistry.Descriptors()` no longer drops the first searched
  registry (an off-by-one in the reverse loop). (#7)
- **sr:** `RegisterByDescriptor` no longer overwrites an explicitly provided
  `Descriptor.Health` with the instance's own `Health`. (#7)
- **hdlm/redislock:** `Unlock` now actually releases the lock (it was a no-op).
  The lock handle is retained so it can be released or refreshed, repeated
  `(Try)Lock` calls refresh the lease instead of reporting the lock as taken, and
  a transient `Release` failure preserves the handle so callers can retry. (#8)
- **hdlm/redislock:** `NewDlm` now initializes the embedded `Health`, preventing a
  nil-interface panic when the DLM is registered as a health check. (#8)
- **hlog:** `SetGlobalLogger` now binds the global `Message` function to the
  driver's `Message` (it was bound to `Info`). (#9)
- **hlog:** `FieldToKeyVal` decodes log fields to their real Go values instead of
  raw integer bits, so bool/float/duration/time fields are no longer reported as
  numbers. (#9)
- **lg:** `UseEmbeddedFieldsInPackage` returns the repackaged embedded fields
  instead of the original slice. (#10)
- **hexa:** `WithBaseTranslator` stores the base translator under the correct
  context key, so `CtxBaseTranslator` works and the translator re-localizes when
  the locale changes. (#11)
- **hexa:** The default context propagator can now extract a context that carries
  no user (it previously required the user key and failed otherwise). (#11)
- **hexa:** User accessors can no longer panic on non-string meta values —
  `validateUserMetaData` validates the string fields up front. (#11)
- **hexa:** `defaultError` implements `Unwrap`, so `errors.Is`/`errors.As`
  traverse to the wrapped internal error. (#12)

### Changed

- **hexa:** After `WithBaseTranslator`, `CtxTranslator` returns the *localized*
  translator and re-localizes on locale change (previously it returned the
  unlocalized base translator). (#11)

### ⚠️ Upgrade notes (observable behavior changes)

- **Stricter user construction:** `NewUserFromMeta` / `MustNewUserFromMeta` /
  `User.SetMeta` now reject meta whose `id`/`email`/`phone`/`name`/`username`
  values are not strings, returning an error at construction (or panicking, for
  the `Must` variant) where previously construction succeeded and only the
  accessor panicked. (#11)
- **`hlog.FieldToKeyVal` output types:** custom log drivers that consume this
  exported function now receive properly-typed values (e.g. `bool`,
  `time.Duration`) instead of `int64`. (#9)
- **Sentry now emits messages:** with the Sentry driver installed, `hlog.Message`
  now triggers `CaptureMessage`; expect message-level events that were previously
  dropped. (#9)
- **`errors.Is`/`errors.As`** against a hexa error now also match its internal
  cause. (#12)

### Compatibility

No exported symbols were removed, renamed, or had their signatures changed.
`hurl.ResponseErr` is newly exported (additive); the `hurl` package did not
compile before #6.

## [0.1.0]

Initial tagged release.

[Unreleased]: https://github.com/Kamva/hexa/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Kamva/hexa/releases/tag/v0.1.0
