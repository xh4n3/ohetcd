# ohetcd
Object Tree Mapping(OTM) for etcd, just like ORM!

## Why OTM?

Because etcd is watchable and distributed by its nature, OTM makes it simpler.

## Scenario

Use this when you need a distributed and consistent variable across many copies of your application.

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

## Tips
ohetcd use JSON to marshal and unmarshal object, so set `Tag` for every field.