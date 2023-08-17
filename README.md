# TL;DR
A wip `devdraw` protocol interposer for debugging

# Rationale
At some point during the course of developing [rjkroege/edwood: Go version of Plan9 Acme Editor](https://github.com/rjkroege/edwood), I thought that it would be nice to have a new implementation of `devdraw`. This seemed
like a large enough project that I wanted to have some incremental milestones. The first such
milestone was an *interposer* `devdraw` that recorded the `devdraw` protocol to a file while delegating the implementation to the original `devdraw`.

