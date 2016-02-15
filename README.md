# Scape-Go

> Un-obstrusive docker helper for go build environments

## TL;DR

Mount your project:

```docker run -t -d -v `pwd`:/src -e MY_REPO=foo/bar --name testing dmp42/scape-go bash```

Now, just do whatever you would do usually to test your code - say `go fmt`, just do:

```docker exec testing my go fmt```

Or just call your makefile, or whatever solution you have in place:

```docker exec testing my make foo bar baz```

... edit your code, rinse, repeat.

## Sugar

### Multiple go versions management

```docker exec testing my gvm list```

By default, the container starts with "future", which currently points at go1.6rc2. Fancy aliases are: `future`, `stable`, `old`, `1.6`, `1.5`, `1.4`.

Use explicit version names if you want to pin your tasks to a specific go environment. Alternatively, use a digest for the scape-go container.

```docker exec testing my gvm use stable```

```docker exec testing my gvm list```

```docker exec testing my go version```

### Package list without the vendor pain

```docker exec testing my echo \$MY_PKGS```

### Build tags

Have some? Just pass them at mount time:

```docker run -t -d -v `pwd`:/src -e MY_REPO=foo/bar -e MY_BUILDTAGS=foo --name testing dmp42/scape-go bash```

### Fast

Yes, it's fast. The overhead against a native go environment is minimal.