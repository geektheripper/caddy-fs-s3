# Caddy FS module for S3

stolen from [sagikazarmark/caddy-fs-s3](https://github.com/sagikazarmark/caddy-fs-s3)

use dsn to specify the s3 bucket, refer to [s3_dsn](https://github.com/geektheripper/vast-dsn/tree/goshujin-sama/s3_dsn)

```shell
xcaddy build --with github.com/geektheripper/caddy-fs-s3
```

## Usage

```caddyfile
{
	filesystem bucket1 s3 {
        dsn "s3://minioadmin:minioadmin@localhost:9000/bucket1?force-path-style=true&protocol=http"
	}
	filesystem bucket2 s3 {
        dsn "$(S3_BUCKET_DSN)"
	}
}

example.com {
    file_server {
        fs bucket1
    }
}
spa.com {
    handle_path /api* {}
    handle {
      fs bucket2
      root web
      encode zstd gzip
      try_files {path} /index.html
      file_server
    }
}
```
