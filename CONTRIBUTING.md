# Contributing to gofeed

## Reporting a problem feed

Most gofeed bugs are "this feed parses wrong" bugs. Open an issue and include
the feed URL or, if the feed may change or disappear, the relevant XML/JSON
snippet. That's usually enough to work with.

## Before opening a pull request

For bug fixes, go ahead and open a PR directly. For new features or behavior
changes, open an issue first so the design can be discussed before you invest
time in an implementation.

PRs target `master`. CI runs the following on the two most recent Go releases;
running them locally first saves a round trip:

```bash
go build ./...
go vet ./...
staticcheck ./...        # go install honnef.co/go/tools/cmd/staticcheck@latest
go test -race -shuffle=on ./...
```

## Test fixtures

Parser behavior is verified with fixture pairs in `testdata/parser/{rss,atom,json,universal}`:
an input file (`name.xml` or `name.json`) and the expected parse result
(`name.json`). The tests glob these directories, so adding a pair is all it
takes — no test code required. Fixes for reported bugs are conventionally named
after the issue, e.g. `issue_217_enclosure_children.xml`.

If your change affects how format-specific fields map to the universal `Feed`
type, the same pattern applies under `testdata/translator/`.

A parser fix or feature without a fixture demonstrating it will be asked to
add one.
