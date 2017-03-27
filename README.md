Gotta Watch Your Code!
======================

[![Build Status](https://travis-ci.org/makii42/gottaw.svg?branch=master)](https://travis-ci.org/makii42/gottaw)

This is yet another funny daemon that watches a file system folder to 
check for change, and execute a command once that happens. 

This thing is still pretty basic and requires a config file right now. Excludes are pretty cumbersome, as globbing in pretty basic in default go (and I did not yet bother looking for something better).

So, check out the [.gottaw.yml](https://github.com/makii42/gottaw/blob/master/.gottaw.yml) here:

    excludes:
      - gottaw
      - .git
      - .vscode
    pipeline: 
      - go build -v .
      - go test -v ./... 
      - go install

But, after this burdon, it works pretty well!

Some of the recent changes:

- [x] Improve Globbing, e.g. .git/**
- [x] Auto-Track addition and removal of folders
- [x] Auto-Reload config if config file is changed
- [x] Growl support
- [x] Pre-define sensible defaults for various setups (go, node, ... )

Obviously it still needs a LOT of polish. I don't have a lot on the agenda 
except optimizations. I'd be especially grateful if someone on Windows 
takes a look on the build and defaults...

PRs Welcome!