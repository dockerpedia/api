## How to use?


#### Search images of the user

```
curl -X POST \
  http://localhost:8080/api/v1/viz \
  -H 'content-type: application/json' \
  -d '{ "user": "google"}'
```

#### Search images with an installed package

```
curl -X POST \
  http://localhost:8080/api/v1/viz \
  -H 'content-type: application/json' \
  -d '{ "package": "nginx"}'
```
