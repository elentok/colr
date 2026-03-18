# CI

Please use github actions to setup CI jobs for running the tests and for
releasing packages.

Use ~/dev/gx/main as reference.

## Tests

- The tests should run on both Mac and Ubuntu

## Releases

- Use goreleaser to build binaries and create a brew package
- Use my brew tap (github.com/elentok/homebrew-stuff), I've already set up
  GORELEASER_HOMEBREW_GITHUB_TOKEN in the GitHub repo settings
