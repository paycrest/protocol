data "remote_dir" "staging" {
  name = "protocol-staging"
}

env "staging" {
  name = atlas.env
  url  = getenv("DATABASE_URL")
  migration {
    dir = data.remote_dir.staging.url
  }
}

data "remote_dir" "production" {
  name = "protocol-production"
}

env "production" {
  name = atlas.env
  url  = getenv("DATABASE_URL")
  migration {
    dir = data.remote_dir.production.url
  }
}
