# Security Model

Status: Initial threat model v0.1

## 1. Authorized use

KryptaSec is designed for systems the operator owns or is explicitly authorized to assess.

Authorization must be enforced by product design rather than only by README text.

## 2. Security boundaries

The initial trust boundaries are:

1. **Main KryptaSec process** — trusted orchestration code.
2. **Agent/LLM output** — untrusted input.
3. **Target application** — untrusted.
4. **Sandbox worker** — disposable, constrained execution.
5. **Browser content** — hostile/untrusted.
6. **Model provider** — external service unless local.
7. **Generated patches** — untrusted until reviewed and tested.

No LLM response is treated as an authorization decision.

## 3. Scope enforcement

Remote active tests require an explicit target policy.

Policy can constrain:

- scheme;
- hostname/domain;
- IP/CIDR;
- port;
- URL path;
- request count;
- concurrency;
- time budget;
- passive vs active mode.

Rules:

- redirects do not extend scope;
- DNS results do not extend scope;
- links/crawled hosts do not extend scope;
- agent instructions do not extend scope;
- discovered subdomains are not automatically authorized;
- an out-of-scope request must fail closed.

## 4. Sandbox

Potentially dangerous commands run only in an isolated runtime.

Baseline requirements:

- non-root process where possible;
- no privileged containers;
- no host Docker socket;
- read-only mounts by default;
- bounded memory/CPU/PIDs;
- execution timeout;
- explicit network egress policy;
- ephemeral filesystem;
- cleanup after cancellation/failure;
- only whitelisted artifacts exported.

## 5. Prompt/tool injection

Source files, web pages, HTTP responses and issue text can contain adversarial instructions.

Rules:

- target content is data, not system policy;
- tools enforce authorization independently of agent prompts;
- privileged operations require typed parameters and policy checks;
- secrets are not made available to arbitrary tools;
- model-generated shell commands never bypass runtime policy.

## 6. Secret handling

KryptaSec must distinguish:

- operator secrets;
- credentials supplied for authorized testing;
- secrets discovered inside target source;
- tokens created by test fixtures.

Requirements:

- redact secrets from logs and reports by default;
- avoid sending discovered secrets to external LLMs;
- never persist provider API keys in scan artifacts;
- encrypt sensitive persistent configuration if local key storage is introduced;
- provide explicit deletion of scan artifacts.

## 7. Findings integrity

A finding has explicit validation status.

```text
hypothesis
candidate
validated
rejected
```

Only `validated` findings may be presented as confirmed exploitation.

Validation evidence must identify:

- tool used;
- relevant request/input;
- relevant response/output;
- timestamp;
- target identity;
- validation constraints.

## 8. Patch safety

Generated patches:

- are written to an isolated worktree;
- are displayed before applying to the user's working branch;
- must not commit/push automatically by default;
- must be build/test checked where supported;
- must be security-retested when tied to a validated finding.

## 9. Denial-of-service controls

Default behavior should minimize load:

- bounded concurrency;
- rate limiting;
- request/time budgets;
- payload size limits;
- no destructive stress testing in normal modes.

Tests whose purpose is availability exhaustion are out of scope for the initial product.

## 10. Logging and privacy

Logs should be structured and support redaction.

Do not log by default:

- Authorization headers;
- cookies;
- API keys;
- full credentials;
- unnecessary response bodies containing personal data.

## 11. Supply-chain security

Before v1.0:

- pin dependencies;
- generate dependency inventory/SBOM;
- run vulnerability scanning;
- review licenses;
- publish checksums;
- sign releases where feasible;
- avoid unverified remote install scripts as the only installation path.

## 12. Vulnerability disclosure

Before public v1.0, publish `SECURITY.md` with:

- supported versions;
- private reporting channel;
- expected response process;
- scope for KryptaSec itself.
