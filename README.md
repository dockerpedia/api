# api

## Search the user's images

```
curl -X POST \
  http://localhost:8080/api/v1/viz \
  -H 'content-type: application/json' \
  -d '{ "user": "google"}'
```

## Search images with a specified package

```
curl -X POST \
  http://localhost:8080/api/v1/viz \
  -H 'content-type: application/json' \
  -d '{ "package": "nginx"}'
```
