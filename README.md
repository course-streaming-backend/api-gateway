# api-gateway

The user-facing HTTP API for the course streaming backend. It is the **only**
service exposed to clients; every backend service is reached over **gRPC** and
has no public HTTP surface.

## Local development with Tilt

Requires `docker`, `tilt`, and `kind`.

```sh
tilt up        # creates a kind cluster, builds images, deploys everything
tilt up auth-service   # or start only one service
```

Tilt provisions the `course-streaming` kind cluster on first run, builds and
hot-reloads the images on file changes, and opens a dashboard at
`http://localhost:10350` with logs, health, and resource status.

- API docs / OpenAPI UI: http://localhost:8080/docs
- Health: http://localhost:8080/api/health

The Tiltfile lives in this repo and references `../auth-service` for its image
and manifests, so `tilt up` must be run from this directory.

## Endpoints

| Method | Path              | Backend                | Status              |
| ------ | ----------------- | ---------------------- | ------------------- |
| GET    | `/api/health`     | – (gateway itself)     | `200`               |
| GET    | `/api/auth/health`| auth-service via gRPC  | `200`               |
| POST   | `/api/auth/register` | auth-service via gRPC| `501` (not implemented) |
| POST   | `/api/auth/login` | auth-service via gRPC  | `501` (not implemented) |

## Proto regeneration

```sh
buf generate
```

