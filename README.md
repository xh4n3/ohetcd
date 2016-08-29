# ohetcd
Object Tree Mapping(OTM) for etcd, just like ORM!

## Why OTM?

Because etcd is watchable and distributed in nature, OTM makes it simpler.

## API

### Set(dir, obj, deep)
set a object obj mapping to a etcd directory dir, if deep is set, save recursively, if not, save it in one node.

### Update()
manually update a object from etcd.

### Save()
save newest object onto etcd.

### Watch()
use etcd built-in watch.

### Unwatch()
stop watching.
