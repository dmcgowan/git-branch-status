# git branch-status

`branch-status` is a git utility to making managing large number of branches
either across many remotes easier. Branch status allows comparing all branches
against their upstream or any arbitrary branch to show the number of 
commit differences. This aids in discovering branches which are out of sync
from their own remotes as well as branches that need to be rebased.

The output of `branch-status` shows a comparison between a branch and either
their upstream or another branch. The comparison shows the number of commits
since the common ancestor in both the target branch (left) and 
upstream branch (right).

## Installation

Install via `go get`

```
$ go get github.com/dmcgowan/git-branch-status
```

Don't have go? Install from [golang.org](http://golang.org/doc/install)

## Usage

Ensure the `$GOPATH` `bin` directory is on `$PATH` or that the
`git-branch-status` binary is on the path. Optionally each command
invoked directly by running `git-branch-status`.

```
$ git branch-status
...
```

### How to interpret output

#### Rebase needed
```
some-branch   4|15    upstream/some-branch
```
This remote branch contains 15 commits
which are not contained in the local branch. Likewise, the local
branch has 4 commits which have been added locally and not pushed
upstream. To rebase, the local branch will need to be rewinded and
the 4 commits replayed on top of the upstream. Run `pull` with the
`--rebase` flag to resolve `git pull --rebase upstream some-branch`.

#### Local out of sync
```
some-branch   0|15    upstream/some-branch
```
The remote branch contains commits which are not part of the local
branch. A simple pull will resolve `git pull upstream some-branch`
by fast forwarding the local branch.

#### Remote out of sync
```
some-branch   3|0     upstream/some-branch
```
The local branch contains changes which have not been pushed upstream.
If the branch is ready to be pushed, then this can be resolved
by pushing the branch `git push upstream some-branch`.

## Common use cases

### List out-of-sync branches
List all branches which are out of sync from their upstream.

```
$ git branch-status
v2-registry-command-header        191|26    stevvooe/v2-registry-command-header
```

Show all branches even those in sync with their upstream.
```
$ git branch-status -a
9468-cert-path                      0|0     jfrazelle/9468-cert-path
add_hostname_docker_info            0|0     vieux/add_hostname_docker_info
address-digest-deadlock             0|0     stevvooe/address-digest-deadlock
auth-option                         0|0     bfirsh/auth-option
concurrent-pull-fix                 0|0     stevvooe/concurrent-pull-fix
manifest-close-archive              0|0     rhvgoyal/manifest-close-archive
master                              0|0     origin/master
parallel-load                       0|0     dougm/parallel-load
registry-info-refactor              0|0     lindenlab/registry-info-refactor
registry_auth_refactor              0|0     jlhawn/registry_auth_refactor
tls_libtrust_auth                   0|0     origin/tls_libtrust_auth
tls_libtrust_auth-documentation     0|0     bfirsh/tls_libtrust_auth-documentation
v2-registry-auth                    0|0     brianbland/v2-registry-auth
v2-registry-command-header        191|26    stevvooe/v2-registry-command-header
v2-registry-tests                   0|0     icecrime/v2-registry-tests
wip_provenance                      0|0     origin/wip_provenance

```

Run `git fetch --all` to ensure all remote upstreams are in sync.

### Check rebases against master

`branch-status` allows comparing branches against any other branch.
This can be useful for checking when a branch was rebased against
another, often useful when needing to keep a branch rebased off of
master. The `-sort` flag can be used to sort by the most out of sync
branches, use `-sort=right` in this case.

Example
```
$ git branch-status -sort=right master
distribution-refactor                   5|10    master
distribution-refactor-tibor             7|10    master
concurrent-pull-fix                     1|59    master
trust-only-pull-by-digest              10|171   master
trust-demo-tuf-push-pull               16|171   master
vendor-distribution                     1|558   master
use-distribution-api                    4|1270  master
fix-official-image-management           1|2206  master
v2-push-test                            2|2494  master
v2-error-handling                       1|2495  master
address-digest-deadlock                 1|2556  master
parallel-load                           1|3201  master
registry-split                          4|3287  master
manifest-close-archive                  2|3325  master
search-refactor-test                    1|3492  master
remote-client-signature                 1|3511  master
v2-registry-basic-auth-fix              1|3610  master
registry-refactor                       1|3674  master
endpoint-creation-refactor              2|3674  master
graph-unit-tests                        4|3674  master
v1-auth-fix                             1|3755  master
v2-registry-mirroring                   1|3755  master
registry-panic-fix                      2|3755  master
v2-registry-fallback                    4|3790  master
v2-registry-command-header              2|3791  master
lk4d4-add_tests_for_registry           22|3980  master
v2-registry-refactor                   25|3980  master
v2-registry-tag-fix                    18|4145  master
stevvooe-v2-registry                    9|4288  master
v2-registry-pushpull                   12|4288  master
push-official-v2                        8|4308  master
registry_auth_refactor                  1|4381  master
9468-cert-path                          1|4443  master
v2-registry-auth                        4|4624  master
provenance                              9|4624  master
add_hostname_docker_info                2|4783  master
libtrust_key                            2|4785  master
network_plugin                          1|4992  master
plugin_proposal                         1|4992  master
tls_libtrust_auth-documentation         6|5198  master
auth-option                            17|5198  master
provenance_resumable                    3|5467  master
registry_v2_switch                      1|5511  master
wip_provenance                         40|6220  master
```

## License

See LICENSE

