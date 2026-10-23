# How our versioning works

> We do **not** use a strict semantic versioning pattern

We use a well known pattern: MAJOR.MINOR.PATCH

MAJOR changes are for versions that break backwards compatibility
MINOR are for new features and bug fixes, while maintaining backwards compatibility
PATCH are for bug fixes only, while maintaining backwards compatibility

The `VERSION` file holds the next (or current) release number.
Git tags use a `v` prefix (`v4.3.6`) and must match `VERSION`.
Run Jenkins on wrapper `edge` first; only then push the tag on GitHub `edge` to start the release workflow (see RELEASING.md).

For some projects, we add a short git-sha to the version string, denoted with a '-' (ex. 3.1.0-ab2d12).
This appendix has no meaning for precedence or does not mean anything regarding changes implemented.