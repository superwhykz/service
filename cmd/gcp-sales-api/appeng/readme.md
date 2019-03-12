# Appengine

Some steps to deploy services... need more descriptions

Thoughts on **standard** vs **flex**

I would really only use flex if I needed the service to host some non-HTTP protocol or _absolutely_ needed Cloud MemoryStore. Otherwise just use GAE with the built in Memcache if needed(a cache that bad).

## Dispatch deploy

Higher level of abstraction to manage routes among the services.

```bash
$ gcloud app deploy cmd/gcp-sales-api/appeng/dispatch.yaml
```

## Deploy services

This step can be automation during image build via Cloud Build service. Below covered manual steps.

### Dev environment

```bash
$ gcloud app deploy cmd/gcp-sales-api/appeng/dev/gcp-sales-api-app.yaml --image-url=gcr.io/PROJECT_ID_NEEDED/gcp-sales-api:0.0.1-dev
```

### Prod environment

```bash
$ gcloud app deploy cmd/gcp-sales-api/appeng/prod/gcp-sales-api-app.yaml --image-url=gcr.io/PROJECT_ID_NEEDED/gcp-sales-api:0.0.1
```
