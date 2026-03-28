# Security Policy

## Supported Scope

Security fixes are prioritized for:

- the active `main` branch
- the most recent tagged release series

This scope includes the Go module, the `gotestout` subpackage, repository automation, and tracked release workflows.

## Reporting A Vulnerability

Use GitHub private vulnerability reporting for this repository if that option is available to you.
If private reporting is unavailable, contact the maintainer privately before opening any public issue or pull request.

Please include:

- affected package or file path
- impact summary
- reproduction steps or proof
- affected versions or commit ranges if known
- any mitigation or patch ideas you already have

## Response Targets

- initial acknowledgment target: 5 business days
- triage target: 10 business days

Remediation timing depends on severity, exploitability, and release timing.

## Disclosure Guidance

- do not publish proof-of-concept exploit details before a fix or coordinated mitigation is available
- maintainers may request coordinated disclosure timing
- security fixes should include tests and release notes when behavior changes

## Out Of Scope

- general feature requests
- non-security correctness bugs
- style or color/theme preferences
- local environment setup problems that do not create a security impact
