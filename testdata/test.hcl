// This uses the dynblock extension to create dynamic blocks
// from an array.
dynamic "server" {
  for_each = bigones
  labels = [ title(server.value) ]
  content {
    addr = urlify(server.value)
  }
}

// Here is an example of using an interpolated template
// function.
server "Beer dot com" {
  addr = "${upper(proto)}://www.beer.com"
}
