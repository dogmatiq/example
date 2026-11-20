# Public demo deployment

This directory configures the public demo at https://example.dogmatiq.io. It is
not needed to write, test, or run the example locally; see the repository README.

Shared CI applies this configuration after publishing and checking a `v*` Docker
image. `RAILWAY_IMAGE_TAG` selects that exact image, and the environment-scoped
`RAILWAY_TOKEN` selects the project. Image polling is disabled for the bank.

To preview a change locally, install Node.js 24 and Railway CLI 5.75.0, then run:

```sh
npm ci --prefix .railway
RAILWAY_IMAGE_TAG=v0.2.1 railway config plan
```

Authenticate with `railway login` and link the intended project/environment, or
provide its `RAILWAY_TOKEN`. Review the plan before applying; CI does not permit
destructive changes automatically.

The public domain is attached only in the `production` environment. An isolated
test environment can deploy the same configuration without claiming that domain.
