# Legal and Provenance Policy

This is an engineering provenance policy, not legal advice.

## 1. Independent implementation

KryptaSec may implement capabilities that exist in other application-security products. Functional goals and publicly known security techniques are not an instruction to copy another project's protected expression.

Changing programming language alone does **not** make copied implementation independent.

## 2. Clean-room rule

Unless a specific third-party component is intentionally incorporated under a compatible license and recorded, contributors must not copy:

- source code;
- prompts;
- tests;
- documentation text;
- screenshots/UI layouts;
- logos, illustrations or branding;
- proprietary datasets;
- internal identifiers or distinctive naming;
- generated patches copied from another product's repository.

Reference products may be studied at the level of public capabilities, protocols, standards and user-visible behavior. KryptaSec implementation should be derived from its own specification and tests.

## 3. Strix reference boundary

Strix is a market/reference product for capability comparison. KryptaSec is not to be implemented as a line-by-line, file-by-file or prompt-by-prompt rewrite.

Do not:

- translate Strix Python files into Go;
- preserve its module/file organization merely with renamed folders;
- copy its prompts and rewrite wording;
- recreate screenshots pixel-for-pixel;
- copy its docs, examples or tests;
- use "same behavior because the source does X" as implementation provenance.

Instead:

1. define the desired capability in KryptaSec docs;
2. define its inputs, outputs and security constraints;
3. implement from that specification using independent code;
4. validate against KryptaSec-owned fixtures.

## 4. Third-party code

When intentionally reusing third-party code:

- record the source repository/package;
- record exact version/commit;
- record license;
- retain required notices;
- document modifications if required;
- ensure compatibility with KryptaSec's eventual distribution license.

A future `THIRD_PARTY_NOTICES.md` will be generated/maintained before public releases.

## 5. Dependency provenance

Every direct dependency should be attributable to:

- package name;
- version;
- official source;
- license;
- reason for use.

Automated license/SBOM tooling should be introduced in Phase 0/1.

## 6. Branding

"KryptaSec" branding, logo, website copy and visual language should be original.

Do not imply affiliation with Strix or any other vendor.

## 7. Distribution license

The repository is currently public, but the final distribution license has not yet been selected.

Before the first release, decide whether KryptaSec will be:

- OSI open source;
- source-available with restrictions;
- dual-licensed;
- proprietary with public source viewing.

This choice must be recorded in an ADR and implemented with a root `LICENSE` file.

## 8. Contribution provenance

Before accepting substantial external contributions, define contributor terms appropriate to the chosen license. Options include Developer Certificate of Origin (DCO) or a Contributor License Agreement (CLA).

Until then, maintainers should avoid merging ambiguous third-party code drops whose provenance cannot be established.
