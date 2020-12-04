# Fetch Artifact

Fetch artifact is a tool for downloading artifacts from Android's continuous integration service.


## Options

* `target`: **Required** - The target you would like to download the artifact from.
* `build_id`: **Required** - The build_id of the target to download the artifact from.
* `artifact`: **Required** - The artifact to download.
* `output`: *Optional* - If you would like the contents of the file to be written to a specific file
* `-`: *Optional* - If you would like the contents of the file to be written to stdout  (must be the last arg)


## Example useage

```
fetch_artifact -target=aosp_arm64-userdebug -build_id=7000390 -artifact=COPIED
```

### Streaming contents to stdout

```
fetch_artifact -target=aosp_arm64-userdebug -build_id=7000390 -artifact=COPIED -
```

## Development

### Building

OUT_DIR=out ./build