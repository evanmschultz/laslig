# Contributing to Läslig

Läslig is a Go library for structured, Charm-native CLI output.
Contributions should keep the library small, readable, and easy to adopt.

## Prerequisites

- Go matching the version in [`go.mod`](./go.mod)
- Mage installed with the same pinned version used in CI:

```bash
go install github.com/magefile/mage@v1.17.0
```

- `vhs` only when you intentionally change tracked demo tapes or README GIFs
- `gh` only when you are maintaining releases or explicitly checking GitHub Actions locally

For routine development, Go and Mage are the only required local tools.

## Daily Workflow

1. Work from the repository root.
2. Keep changes focused and keep docs/tests aligned with behavior.
3. Use the normal local gate before you hand work back:

```bash
mage check
```

4. Use faster or narrower commands while iterating:

```bash
mage test
mage build
mage demo
go test ./...
```

In this repository's bare-root worktree layout, prefer the Mage targets for
local build and demo commands. They already pass `-buildvcs=false` where
needed so local worktree builds do not fail on VCS stamping.

5. Run VHS only when user-visible terminal output or tracked tapes change:

```bash
mage vhs
```

## Golden Snapshots

Structural output is protected by Charm `x/exp/golden` snapshots.
Update them intentionally, not accidentally:

```bash
go test ./internal/examples -args -update
go test ./examples/all -args -update
go test ./examples/gotestout -args -update
go test ./gotestout -run 'TestRenderPlainCompactGolden|TestRenderHumanStyledCompactGolden' -args -update
```

If snapshots change because output changed on purpose, update the README and VHS assets in the same change when relevant.

## Dependency Policy

- Keep runtime/library dependencies narrow: standard library first, Charm packages only when they materially help.
- Do not add `fang` as a core dependency.
- Do not add `charm/log` as a core dependency.
- Demo-only or tooling-only dependencies are acceptable when clearly scoped and documented.
- Prefer patch and minor updates by default.
- Treat major dependency upgrades as design changes that require explicit review.
- Before upgrading Charm-family libraries, Glamour, or `x/*` dependencies, review upstream release notes and compatibility implications.
- After dependency bumps, run:

```bash
go mod tidy
mage check
```

- If rendering changed intentionally, refresh goldens and VHS assets too.

## Documentation And Tests

- Keep public doc comments current when exported behavior changes.
- Keep [`README.md`](./README.md), examples, and tracked VHS assets aligned with shipped behavior.
- Maintain at least 70% statement coverage in every package.
- Prefer real behavior checks over mock-heavy tests.

## Release Notes For Contributors

Normal contributors do not need GoReleaser or GitHub release tooling.
Maintainers handle tags, publishing, and post-publish verification.

## Security

Do not open public issues for suspected vulnerabilities.
Follow [`SECURITY.md`](./SECURITY.md).

## License

By submitting a contribution, you agree that your contribution may be distributed under the terms of the Apache 2.0 license in [`LICENSE`](./LICENSE).
