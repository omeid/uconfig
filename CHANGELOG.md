# Changelog

## v0.9.0

### Added
- **Map support in flat package.** Maps can now be set from flat sources like flags and env vars using `key:value,key:value` syntax (e.g. `-my-map "a:1,b:2"`). Supports all combinations of string, int, uint, and float key/value types, plus `TextUnmarshaler` and `time.Duration` values.

## v0.8.0

### Added
- **Generic API.** `uconfig.Classic[T]()` and `uconfig.Load[T]()` for type-safe configuration loading.
- **Command support.** Last positional argument is treated as a command via `flag:",command"` tag.
- **`uconfig` struct tag.** Rename fields for flags and env vars with dot notation (e.g. `uconfig:"Pre.Sync"` gives `-pre-sync` flag and `PRE_SYNC` env var).

## v0.7

### Breaking Changes

Field names now can be renamed at struct level instead of globally by prefixing them with a `.`.

So that in the following example, the config Redis.Host will be mapped to Redis.Address
While Port will be mapped to REDIS_SERVICE_PORT.


```go
type Redis struct {
  Host string `unconfig:".Address"`
  Port string `unconfig:"REDIS_SERVICE_PORT"`
}


type Config struct {
  Redis Redis
}
```




## v0.6

Artificial release to retract v1.


## v0.5

### Breaking Changes


#### Classic

Classic now automatically prints usage and exists if user passes `-h` or `--help`.



## v0.4

### Breaking Changes

#### plugins/file.Files

And so `uconfig.Files`

Has changed from:
```go
type Files []struct {
	Path      string
	Unmarshal Unmarshal
}
```

to

```go
type Files []struct {
	Path      string
	Unmarshal Unmarshal
	Optional  bool /* new field! */
}
```


#### plugins/files.Plugins

No longer accepts an option.



#### plugins/flag

No longer supports `SetUsage`.
