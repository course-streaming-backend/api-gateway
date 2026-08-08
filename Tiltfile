# Tiltfile - Course streaming backend local dev.
# Run `tilt up` from the api-gateway directory. Requires docker + kind.

cluster_name = 'course-streaming'
kube_context = 'kind-' + cluster_name

# --- Cluster ---------------------------------------------------------------

# Create the kind cluster on first run; this is a no-op afterwards.
existing = str(local('kind get clusters 2>/dev/null || true')).strip()
if cluster_name not in [c.strip() for c in existing.splitlines()]:
    local('kind create cluster --name ' + cluster_name + ' --wait 120s')

allow_k8s_contexts(kube_context)

# --- Images ----------------------------------------------------------------

docker_build('api-gateway', '.', dockerfile='Dockerfile')
docker_build('auth-service', '../auth-service', dockerfile='../auth-service/Dockerfile')
docker_build('auth-migrations', '../auth-service', dockerfile='../auth-service/deploy/migrate/Dockerfile')

# --- Manifests -------------------------------------------------------------

k8s_yaml([
    'k8s/api-gateway.yaml',
    '../auth-service/k8s/auth-service.yaml',
    '../auth-service/k8s/postgres.yaml',
    '../auth-service/k8s/migrations.yaml',
])

# --- Resources -------------------------------------------------------------

k8s_resource(
    'api-gateway',
    port_forwards='8080:8080',
    links=['http://localhost:8080/docs'],
    labels=['gateway'],
)

k8s_resource(
    'auth-service',
    port_forwards='50051:50051',
    resource_deps=['postgres'],
    labels=['services'],
)

k8s_resource(
    'postgres',
    port_forwards='5433:5432',
    labels=['infra'],
)

k8s_resource(
    'migrations',
    resource_deps=['postgres'],
    labels=['infra'],
)
