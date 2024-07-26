# Fetch Artifact

Fetch artifact is a tool for downloading artifacts from Android's continuous integration service.


## Options

* `target`: **Required** - The target you would like to download the artifact from.
* `artifact`: **Required** - The artifact to download.
* **Required**: either `build_id` or `branch`, but not both
  * When only `build_id` is provided, the script would download the artifact from that `build_id`.
  * When only `branch` is provided, the script would download the artifact from the last known good build of that `branch`.
* `output`: *Optional* - If you would like the contents of the file to be written to a specific file.
* `-`: *Optional* - If you would like the contents of the file to be written to stdout  (must be the last arg)


## Example useage

```
fetch_artifact -target=aosp_arm64-userdebug -build_id=7000390 -artifact=COPIED
```

### Streaming contents to stdout

```
fetch_artifact -target=aosp_arm64-userdebug -build_id=7000390 -artifact=COPIED -
```

### Get the latest successful build's artifact without specifying a build_id
```
fetch_artifact -target=aosp_arm64-trunk_staging-userdebug -branch=aosp-main -artifact=COPIED
```

## Development

### Building

OUT_DIR=out ./build