# API Documentation

Base URL: `/api/v1`

Protected base: `/api/v1/protected`

All responses follow the same shape:

```json
{
  "status_code": 200,
  "message": "Success",
  "data": {}
}
```

Paginated responses include a `meta` object:

```json
{
  "status_code": 200,
  "message": "Success",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3,
    "has_next": true,
    "has_previous": false,
    "timestamp": 1700000000
  }
}
```

Error responses use the same wrapper. Validation errors return a list of fields and tags:

```json
{
  "status_code": 400,
  "message": "Bad Request",
  "data": {
    "error": [
      { "field": "Email", "tag": "email" }
    ],
    "meta": {
      "method": "POST",
      "path": "/api/v1/auth/register",
      "endpoint": "/auth/register",
      "status": 400,
      "latency": "2.1ms",
      "ip": "127.0.0.1"
    }
  }
}
```

## Authentication

Protected endpoints require:

```
Authorization: Bearer <access_token>
```

Auth endpoints are throttled (5 requests per minute). Refresh tokens are issued via the `refresh_token` httpOnly cookie.

### POST /auth/register
Request body:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "avatar_url": "https://example.com/avatar.png",
  "password": "secret123"
}
```

Response: `201 Created`

### POST /auth/verify?token=...
Validates the email verification token and sets a refresh token cookie.

Response: `200 OK`
Note: The current implementation echoes the verification token in `data.access_token`.

### POST /auth/login
Request body:

```json
{
  "email": "john@example.com",
  "password": "secret123"
}
```

Response: `200 OK` with `access_token` in `data` and a refresh token cookie.

### POST /auth/refresh-token
Uses the refresh token cookie to mint a new access token.

Response: `200 OK` with `access_token` in `data` and a new refresh token cookie.

### DELETE /auth/logout
Clears the refresh token cookie.

Response: `200 OK`
Note: The current handler returns HTTP 400 while the response body uses `status_code: 200`.

## Users (Protected)

### GET /protected/users/me
Returns the current user profile.

Response: `200 OK`

### PATCH /protected/users/me
Request body:

```json
{
  "name": "John Doe",
  "avatar_url": "https://example.com/avatar.png"
}
```

Response: `200 OK`

### GET /protected/users
Admin-only. Supports pagination with `page` and `limit` query parameters.

Response: `200 OK` (paginated)

### GET /protected/users/:id
Admin-only. `:id` must be a UUID.

Response: `200 OK`
