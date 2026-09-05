data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/atlasloader"
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  dev = getenv("ATLAS_DEV_URL")
  url = getenv("ATLAS_DB_URL")

  migration {
    dir = "file://internal/db/migrations"
  }
}
