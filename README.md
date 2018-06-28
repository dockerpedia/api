# api


```
➜  ~ curl https://api.mosorio.me/api/v1/viz\?query\=redis | prettyjson  | less
```

```
{
    "children": [
        {
            "affilation": null,
            "children": [
                {
                    "analysed": true,
                    "full_size": 7465592,
                    "id": 4671767,
                    "id_docker": 8197954,
                    "image_id": 1366564,
                    "last_check": "2018-04-04T21:35:20.654645Z",
                    "last_try": "2018-04-04T21:35:20.673584Z",
                    "last_updated": "2017-05-13T18:21:58.021674Z",
                    "name": "3.2.8-alpine",
                    "packages": 12,
                    "status": true,
                    "url": "",
                    "username": "",
                    "value": 964,
                    "vulnerabilities_critical": 0,
                    "vulnerabilities_defcon1": 0,
                    "vulnerabilities_high": 2,
                    "vulnerabilities_low": 0,
                    "vulnerabilities_medium": 5,
                    "vulnerabilities_negligible": 0,
                    "vulnerabilities_unknown": 0
                }
                ],
                "description": "Redis is an open source key-value store that functions as a data structure server.",
                "full_name": "library/redis",
                "full_size": 10588033,
                "id": 1366564,
                "last_updated": "2018-04-02T16:28:41.9692Z",
                "name": "redis",
                "namespace": "library",
                "official": null,
                "pull_count": 693183603,
                "repository_type": "image",
                "start_count": 4981,
                "status": true,
                "tags_checked": null,
                "user": "library",
                "value": 1000
        }
    ],
    "name": "redis"
}                
```
