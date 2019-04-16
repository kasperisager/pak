workflow "Integrate" {
  on = "push"
  resolves = ["test"]
}

action "build" {
  uses = "docker://golang"
  runs = "go"
  args = "build ./..."
}

action "test" {
  uses = "docker://golang"
  needs = ["build"]
  runs = "go"
  args = "test ./..."
}
