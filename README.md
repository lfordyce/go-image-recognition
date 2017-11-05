# go-tensorflow-docker

---
- To run execute:
```
$ docker-compose -f docker-compose.yaml up -d --build
```

- Post with curl command:
```
$ curl localhost:8080/recognize -F 'image=@./dog.jpg'
```
- Or with `httpie`
```
http -f POST localhost:8080/recognize image@dog4.jpg
```
- Response:
```json
HTTP/1.1 200 OK
Content-Length: 284
Content-Type: application/json
Date: Wed, 25 Oct 2017 22:15:28 GMT

{
    "filename": "dog4.jpg",
    "labels": [
        {
            "label": "golden retriever",
            "probability": 0.8646044
        },
        {
            "label": "Labrador retriever",
            "probability": 0.11315775
        },
        {
            "label": "beagle",
            "probability": 0.006965436
        },
        {
            "label": "kuvasz",
            "probability": 0.0033914202
        },
        {
            "label": "tennis ball",
            "probability": 0.00308104
        }
    ]
}
```

- Post multiple images to the api:
```
curl localhost:8080/series -F 'images=@./dog.jpg -F 'images=@./dog2.jpg -F 'images=@./dog3.jpg -F 'images=@./dog4.jpg
```