# GoChat



### ERD

![img-erd](https://github.com/user-attachments/assets/6c355370-596c-48a9-84ef-f0c47bdba5e4)

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

### How to Run

