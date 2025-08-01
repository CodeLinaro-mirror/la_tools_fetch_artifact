# Fetch Artifact

Fetch artifact is a tool for downloading artifacts from Android's continuous integration service.


## Options

* `target`: **Required** - The target you would like to download the artifact from.
* `artifact`: **Required** - The artifact to download.
* **Required**: either `build_id` or `branch`, but not both
  * When only `build_id` is provided, the script would download the artifact from that `build_id`.
  * When only `branch` is provided, the script would download the artifact from the last known good build of that `branch`.
* `output`: *Optional* - If you would like the contents of the file to be written to a specific file.
* `client_id`: *Optional* - If authorization is required to download the artifact, please set this parameter as your OAuth2.0 Client ID.
* `secret`: *Optional* - If authorization is required to download the artifact, please set this parameter as your OAuth2.0 Client secret.
* `port`: *Optional* - If you would like to specify the OAuth callback port to listen on. Default: 10502
* `project_id`: *Optional* - The project id being used to access the fetch APIs.
* `-`: *Optional* - If you would like the contents of the file to be written to stdout  (must be the last arg)

## Example usage

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

### Using OAuth to fetch restricted artifacts
In this case, you might need to create an OAuth 2.0 Client ID for a web application and set the redirect URI to `http://localhost:<port>`(default port: 10502).

```
fetch_artifact -target=<restricted_target> -build_id=<id> -artifact=COPIED -client_id=<OAuth_client_id> -secret=<OAuth_client_secret>
```

If you are accessing the fetch APIs from a different project than your OAuth client, you will need to specify the `-project_id` flag:

```
fetch_artifact -target=<restricted_target> -build_id=<id> -artifact=COPIED -client_id=<OAuth_client_id> -secret=<OAuth_client_secret> -project_id=<project_id>
```

## Development

### Building

OUT_DIR=out ./build