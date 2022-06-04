# Changelog

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
