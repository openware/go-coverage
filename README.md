# Go Coverage

This tool is used to export summary of XML coverage report into a SQL database.

## Usage

```bash
export ANALYTICS_DATABASE_HOST="localhost"
export ANALYTICS_DATABASE_PORT="5433"
export ANALYTICS_DATABASE_USERNAME="postgres"
export ANALYTICS_DATABASE_PASSWORD=""
export ANALYTICS_DATABASE_NAME="analytics"
export DRONE_REPO="openware/go-coverage"
export ANALYTICS_COMPONENT="go-coverage"        # use it to differenciate different applications in a mono-repo.
export DRONE_TAG="1.0.0"                        # trigger this script on drone tag to track only stable versions

go run ./ coverage.xml
```

## Postgresql

To create a Postgres instance with the `coverage_analytics` database by Helm

### Install postgresql
```sh
helm upgrade -i postgres-coverage bitnami/postgresql -f ./values.yml -n core
```

values.yml
```yaml
image:
  tag: 14.5.0

auth:
  database: coverage_analytics
  enablePostgresUser: true
  existingSecret: ""
  password: changeme
  postgresPassword: changeme
  replicationPassword: ""
  replicationUsername: repl_user
  secretKeys:
    adminPasswordKey: postgres-password
    replicationPasswordKey: replication-password
    userPasswordKey: password
  usePasswordFiles: false
  username: coverage_analytics
```

### Expose postgresql service (optional)

```sh
kubectl apply -f ./service.yml -n ${your-namespace}
kubectl apply -f ./ingress.yml -n ${your-namespace}
```

service.yml
```yml
apiVersion: v1
kind: Service
metadata:
  name: postgres-coverage-postgresql
  labels:
    app.kubernetes.io/component: primary
    app.kubernetes.io/instance: postgres-coverage
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: postgresql
    helm.sh/chart: postgresql-12.1.2
  annotations:
    meta.helm.sh/release-name: postgres-coverage
    meta.helm.sh/release-namespace: core
spec:
  ports:
    - protocol: TCP
      port: 5432
      targetPort: 5432
      nodePort: 30432
  selector:
    app.kubernetes.io/component: primary
    app.kubernetes.io/instance: postgres-coverage
    app.kubernetes.io/name: postgresql
  type: NodePort
  sessionAffinity: None
  externalTrafficPolicy: Cluster
```

ingress.yml
```yml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: postgres-coverage
  labels:
    app: postgres-coverage
  annotations:
    kubernetes.io/ingress.class: nginx
    kubernetes.io/tls-acme: 'true'
spec:
  tls:
    - hosts:
        - pg.example.app
      secretName: postgres-coverage-tls
  rules:
    - host: pg.example.app
      http:
        paths:
          - path: /
            backend:
              serviceName: postgres-coverage-postgresql
              servicePort: 5432
```
