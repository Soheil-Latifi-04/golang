# go-CRM

A minimal CRM API written in Go. It manages leads (name, company, email, phone) through a small REST interface backed by SQLite.

## Stack

- [Fiber v3](https://github.com/gofiber/fiber) — HTTP router
- [GORM](https://gorm.io/) with the SQLite driver — persistence
- SQLite (`test.db`) — storage, created automatically on first run

## Getting started

```bash
git clone https://github.com/Soheil-Latifi-04/golang.git
cd golang/go-CRM
go mod tidy
go run main.go
```

The server starts on `http://localhost:3000` and creates `test.db` in the working directory on first launch.

## API

| Method | Endpoint              | Description         |
|--------|------------------------|----------------------|
| GET    | `/api/v1/leads`         | List all leads       |
| GET    | `/api/v1/lead/:id`      | Get a single lead    |
| POST   | `/api/v1/lead`          | Create a new lead     |
| DELETE | `/api/v1/lead/:id`      | Delete a lead         |

### Lead object

```json
{
  "Name": "Jane Doe",
  "Company": "Acme Inc",
  "Email": "jane@acme.com",
  "Phone": "555-0100"
}
```

## Example session

With the server running (`go run main.go`), here's a full round trip against a fresh database.

**Create a lead**

```bash
curl -X POST http://localhost:3000/api/v1/lead \
  -H "Content-Type: application/json" \
  -d '{"Name":"Jane Doe","Company":"Acme Inc","Email":"jane@acme.com","Phone":"555-0100"}'
```

```json
{
  "ID": 1,
  "CreatedAt": "2026-08-18T10:15:32Z",
  "UpdatedAt": "2026-08-18T10:15:32Z",
  "DeletedAt": null,
  "Name": "Jane Doe",
  "Company": "Acme Inc",
  "Email": "jane@acme.com",
  "Phone": "555-0100"
}
```

**List all leads**

```bash
curl http://localhost:3000/api/v1/leads
```

```json
[
  {
    "ID": 1,
    "CreatedAt": "2026-08-18T10:15:32Z",
    "UpdatedAt": "2026-08-18T10:15:32Z",
    "DeletedAt": null,
    "Name": "Jane Doe",
    "Company": "Acme Inc",
    "Email": "jane@acme.com",
    "Phone": "555-0100"
  }
]
```

**Get a single lead**

```bash
curl http://localhost:3000/api/v1/lead/1
```

```json
{
  "ID": 1,
  "CreatedAt": "2026-08-18T10:15:32Z",
  "UpdatedAt": "2026-08-18T10:15:32Z",
  "DeletedAt": null,
  "Name": "Jane Doe",
  "Company": "Acme Inc",
  "Email": "jane@acme.com",
  "Phone": "555-0100"
}
```

**Delete a lead**

```bash
curl -i -X DELETE http://localhost:3000/api/v1/lead/1
```

```
HTTP/1.1 204 No Content
```

**Lead not found**

```bash
curl -i http://localhost:3000/api/v1/lead/999
```

```
HTTP/1.1 404 Not Found
{"error":"Lead not found"}
```

## Project structure

```
go-CRM/
├── main.go              # entrypoint, routes, DB setup
├── database/
│   └── database.go       # shared DB connection
└── lead/
    └── lead.go            # Lead model and handlers
```

## Notes

This is a learning/portfolio project focused on a clean Fiber + GORM setup rather than production hardening — there's no auth, validation is minimal, and it uses a local SQLite file rather than a managed database.
