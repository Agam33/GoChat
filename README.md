# GoChat

The goal of this project is to demonstrate the core concepts of a chat system, where users can create rooms and communicate with other users in real time.

### ERD

![img-erd](https://github.com/user-attachments/assets/2c9e9ed6-b84f-4bb0-a612-4b9bde65d3e4)

---

### Tech Stack
- **Language:** Go
- **Web Framework:** Gin
- **ORM:** Gorm
- **Database:** PostgreSQL
- **WebSocket:** [Gorilla](https://github.com/gorilla/websocket)
- **Auth:** JWT
- **Configuration:** Viper
- **Message Queue:** RabbitMQ 
---

### API Documentation

#### REST

This project used [Swagger](https://github.com/swaggo/swag) for REST api documentation. You can access it at: http://localhost:8000/swagger/index.html#/

![img-swagger](https://github.com/user-attachments/assets/5dfe5fc0-3d62-419e-a2fb-445c4ef9c06d)


#### Websocket

##### Connect with [wscat](https://github.com/websockets/wscat)

```bash
wscat -c ws://localhost:8000/api/ws -H "Authorization: Bearer <token>"
```
---

##### Websocket events

[ws_handler.go](https://github.com/Agam33/GoChat/blob/dev/internal/websocket/ws_handler.go) where the events are handled.

`base event msg`
```json
{
    "action": "",
    "data": {} 
}
```

`action` **room_join**
```json
{
    "roomId": 0
}
```

`action` **room_leave**
```js
{
    "roomId": 0
}
```

`action` **room_send_text**
```js
{
    "roomId": 0,
    "text": ""
}
```

`action` **room_reply_text**
```js
{
    "roomId": 0,
    "senderId": 0,
    "replyTo": 0, // msg id
    "text": ""
}
```

`action` **room_delete_message**
```js
{
    "roomId": 0,
    "senderId": 0,
    "messageId": 0
}
```

---

### Response in room:
`/room/:id/messages`
```json
{
    "message": "success",
    "data": [
        {
            "id": 1768993616933,
            "sender": {
                "id": 1768287960356852,
                "name": "agam",
                "imgUrl": null
            },
            "content": {
                "text": "Gimana kabarnya Riswan? masih di Jakarta?",
                "contentType": "text"
            },
            "replyContent": {
                "id": 1768993200764,
                "content": {
                    "text": "halo gam :)",
                    "contentType": "text"
                }
            },
            "createdAt": "2026-01-21T18:06:56.934014+07:00",
            "updatedAt": "2026-01-21T18:06:56.934014+07:00"
        },
        {
            "id": 1768993200764,
            "sender": {
                "id": 1768287970851522,
                "name": "riswan",
                "imgUrl": null
            },
            "content": {
                "text": "halo gam :)",
                "contentType": "text"
            },
            "replyContent": {
                "id": 1768993098,
                "content": {
                    "text": "Halo Riswan",
                    "contentType": "text"
                }
            },
            "createdAt": "2026-01-21T18:00:00.764982+07:00",
            "updatedAt": "2026-01-21T18:00:00.764982+07:00"
        },
        {
            "id": 1768993098,
            "sender": {
                "id": 1768287960356852,
                "name": "agam",
                "imgUrl": null
            },
            "content": {
                "text": "This message was deleted",
                "contentType": "text"
            },
            "replyContent": null,
            "createdAt": "2026-01-21T17:58:18.791927+07:00",
            "updatedAt": "2026-01-21T18:14:00.201961+07:00"
        }
    ],
    "meta": {
        "currPage": 1,
        "nextPage": 2,
        "prevPage": 0
    }
}
```

---

### How to Run
#### Set environment variables
    
See the example `.env` file: [click](https://github.com/Agam33/GoChat/blob/dev/.env.example)

#### Run with Docker
```bash
docker compose up
```
---
