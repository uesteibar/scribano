workflow "CI" {
  on = "push"
  resolves = ["test"]
}

action "test" {
  uses = "docker://docker/compose:1.23.2"
  args = "docker-compose up -d"
}
