import { defineRailway, image, postgres, project, service } from "railway/iac";

// This config is for the public demo. Local setup is documented in ../README.md.
export default defineRailway((ctx) => {
  const tag = process.env.RAILWAY_IMAGE_TAG;
  if (!tag) {
    throw new Error("RAILWAY_IMAGE_TAG must identify a published Docker image");
  }

  const database = postgres("Postgres", { region: "sfo" });
  database.deploy = {
    ...database.deploy,
    limitOverride: {
      containers: { cpu: 1, memoryBytes: 1_000_000_000 },
    },
  };

  const bank = service("bank", {
    source: image(`ghcr.io/dogmatiq/example:${tag}`, {
      autoUpdates: { type: "disabled" },
    }),
    replicas: { sfo: 1 },
    healthcheck: "/",
    healthcheckTimeout: 60,
    deploy: {
      sleepApplication: true,
      restartPolicyType: "ON_FAILURE",
      restartPolicyMaxRetries: 3,
      drainingSeconds: 10,
      overlapSeconds: 0,
      limitOverride: {
        containers: { cpu: 1, memoryBytes: 500_000_000 },
      },
    },
    env: {
      DOGMA_EXAMPLE_ENGINE_DSN: database.env.DATABASE_URL,
      DOGMA_EXAMPLE_PROJECTION_DSN: database.env.DATABASE_URL,
      DOGMA_EXAMPLE_IDLE_TIMEOUT: "1m",
      DOGMA_EXAMPLE_LISTEN_ADDRESS: "[::]:8080",
      PORT: "8080",
    },
    domains: ctx.isEnvironment("production")
      ? [{ domain: "example.dogmatiq.io", port: 8080 }]
      : [],
  });

  // The environment-scoped deployment token selects the project.
  return project(ctx.projectName!, { resources: [database, bank] });
});
