build {
  sources = [
    "source.digitalocean.base_image"
  ]

  provisioner "shell" {
    inline = ["sleep 5"]
  }
}
