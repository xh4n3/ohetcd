# ohetcd
ORM for etcd.

## Set(dir, obj, deep)
set a object obj mapping to a etcd directory dir, if deep is set, save recursively, if not, save it in one node.

## Update()
manually update a object from etcd.

## Save()
save newest object onto etcd.

## Link()
use a goroutine to recursively query for updates.

## Unlink()
stop the goroutine.

## Watch()
use etcd built-in watch.

## Unwatch()
stop watching.
